package uvc

type Branch struct {
	name         string
	headRevision string
	revisions    []Revision
}
