package domain

type NotificationService interface {
	NotifyMarkAdded(mark *Mark)
	// Add other notifications as needed
}
