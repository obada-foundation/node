package search

import (
	"github.com/obada-foundation/node/business/tests"
	"testing"
)

type ServiceTests struct {
	service *Service
}

func TestService(t *testing.T) {
	test := tests.NewIntegration(t)
	t.Cleanup(test.Teardown)

	tests.CreateOwnerObits(t, test)

	service := NewService(test.Logger, test.DB)

	ts := ServiceTests{
		service: service,
	}

	t.Run("search", ts.search)
}

type searchTestCases struct {
	args searchTestCasesArgs
	want searchTestCasesWant
}

type searchTestCasesArgs struct {
	term   string
	offset uint
}

type searchTestCasesWant struct {
	count       int
	total       uint
	lastPage    uint
	currentPage uint
}

func (os ServiceTests) search(t *testing.T) {
	testCases := []searchTestCases{
		{
			args: searchTestCasesArgs{
				term:   "usn1",
				offset: 0,
			},
			want: searchTestCasesWant{
				count:       1,
				total:       1,
				lastPage:    1,
				currentPage: 1,
			},
		},
		{
			args: searchTestCasesArgs{
				term:   "unexists",
				offset: 0,
			},
			want: searchTestCasesWant{
				count:       0,
				total:       0,
				lastPage:    1,
				currentPage: 1,
			},
		},
		{
			args: searchTestCasesArgs{
				term:   "unexists",
				offset: 10,
			},
			want: searchTestCasesWant{
				count:       0,
				total:       0,
				lastPage:    1,
				currentPage: 11,
			},
		},
		{
			args: searchTestCasesArgs{
				term:   `did:obada:owner:678910`,
				offset: 0,
			},
			want: searchTestCasesWant{
				count:       50,
				total:       150,
				lastPage:    3,
				currentPage: 1,
			},
		},
		{
			args: searchTestCasesArgs{
				term:   `did:obada:owner:678910`,
				offset: 1,
			},
			want: searchTestCasesWant{
				count:       50,
				total:       150,
				lastPage:    3,
				currentPage: 2,
			},
		},
		{
			args: searchTestCasesArgs{
				term:   "",
				offset: 0,
			},
			want: searchTestCasesWant{
				count:       50,
				total:       150,
				lastPage:    3,
				currentPage: 1,
			},
		},
		{
			args: searchTestCasesArgs{
				term:   ";",
				offset: 0,
			},
			want: searchTestCasesWant{
				count:       50,
				total:       150,
				lastPage:    3,
				currentPage: 1,
			},
		},
		{
			args: searchTestCasesArgs{
				term:   "'",
				offset: 0,
			},
			want: searchTestCasesWant{
				count:       50,
				total:       150,
				lastPage:    3,
				currentPage: 1,
			},
		},
	}

	for _, tc := range testCases {
		got, err := os.service.Search(tc.args.term, tc.args.offset)

		if err != nil {
			t.Fatalf("service.Search(%q, %d) Error = %s", tc.args.term, tc.args.offset, err)
		}

		if len(got.Obits) != tc.want.count {
			t.Errorf(
				"service.Search(%q, %d) count(Obits) = %d want %d",
				tc.args.term, tc.args.offset, len(got.Obits), tc.want.count,
			)
		}

		if got.Meta.Total != tc.want.total {
			t.Errorf(
				"service.Search(%q, %d) Meta.Total = %d want %d",
				tc.args.term, tc.args.offset, got.Meta.Total, tc.want.total,
			)
		}

		if got.Meta.LastPage != tc.want.lastPage {
			t.Errorf(
				"service.Search(%q, %d) Meta.LastPage = %d want %d",
				tc.args.term, tc.args.offset, got.Meta.Total, tc.want.total,
			)
		}

		if got.Meta.CurrentPage != tc.want.currentPage {
			t.Errorf(
				"service.Search(%q, %d) Meta.CureentPage = %d want %d",
				tc.args.term, tc.args.offset, got.Meta.CurrentPage, tc.want.currentPage,
			)
		}
	}
}
