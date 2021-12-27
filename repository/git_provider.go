package repository

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

const suffix = ".gitignore"

type (
	GitProvider struct {
		tree func() (*object.Tree, error)
	}

	LazyTree  func() (*object.Tree, error)
	Predicate func(term string, file *object.File) bool
	Visitor   func(file *object.File) error
)

// NewGitProvider creates a new provider that uses the gitignore repository.
func NewGitProvider() *GitProvider {
	return &GitProvider{
		tree: initTree(),
	}
}

// Get returns the .gitignore file with the matching name.
func (p *GitProvider) Get(term string) (io.Reader, error) {
	var output io.Reader = nil

	err := p.visit(term, exactPredicate, func(file *object.File) error {
		reader, err := file.Blob.Reader()

		if err != nil {
			return err
		}

		output = reader
		return err
	}, true)

	return output, err
}

// Search searches for files names in the git repository
func (p *GitProvider) Search(term string) ([]string, error) {
	var names []string

	err := p.visit(term, matchPredicate, func(file *object.File) error {
		names = append(names, strings.Replace(file.Name, suffix, "", 1))
		return nil
	}, false)

	if err != nil {
		return nil, err
	}

	return names, err
}

// visit will walk the file iterator
func (p *GitProvider) visit(term string, predicate Predicate, visitor Visitor, followSymbolicLinks bool) error {
	tree, err := p.tree()

	if err != nil {
		return err
	}

	err = tree.Files().ForEach(func(f *object.File) error {
		if predicate(term, f) {
			if followSymbolicLinks && f.Mode == filemode.Symlink {
				r, err := f.Reader()

				if err != nil {
					return err
				}

				buf := new(strings.Builder)
				_, err = io.Copy(buf, r)

				if err != nil {
					return err
				}

				symbolicLink := strings.Replace(buf.String(), suffix, "", 1)

				return p.visit(symbolicLink, predicate, visitor, false)
			}

			return visitor(f)
		}

		return nil
	})

	return err
}

// initTree returns the default git tree for the main branch in the gitignore repo.
func initTree() LazyTree {
	var tree *object.Tree
	var once sync.Once
	var err error

	return func() (*object.Tree, error) {
		once.Do(func() {
			fs := memfs.New()
			storage := memory.NewStorage()

			// We clone the repo into an in-memory store
			r, err := git.Clone(storage, fs, &git.CloneOptions{
				URL:           "https://github.com/github/gitignore",
				SingleBranch:  true,
				ReferenceName: plumbing.ReferenceName("refs/heads/main"),
			})

			if err != nil {
				return
			}

			// To get the tree, we need to get the HEAD commit and then pull the tree.
			// Our purpose is to get the tree for the lastest main branch commit.
			ref, err := r.Head()
			if err != nil {
				return
			}

			commit, err := r.CommitObject(ref.Hash())
			if err != nil {
				return
			}

			tree, err = commit.Tree()
			if err != nil {
				return
			}
		})

		return tree, err
	}
}

// exactPredicate matches a term with the filename (ignoring case)
func exactPredicate(term string, file *object.File) bool {
	lowerTerm := strings.ToLower(fmt.Sprintf("%s%s", term, suffix))
	lowerFilename := strings.ToLower(file.Name)
	return lowerTerm == lowerFilename
}

// matchPredicate preforms a case insensitive regex match
func matchPredicate(term string, file *object.File) bool {
	s := fmt.Sprintf("(?i)%s", term)
	r, _ := regexp.Compile(s)
	return r.MatchString(file.Name)
}
