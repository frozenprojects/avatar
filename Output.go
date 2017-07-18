package avatar

// Output represents a system that saves an avatar locally (in database or as a file, e.g.)
type Output interface {
	SaveAvatar(*Avatar) error
}
