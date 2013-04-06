// Package drive provides access to the Drive API.
//
// See https://developers.google.com/drive/
//
// Usage example:
//
//   import "code.google.com/p/google-api-go-client/drive/v2"
//   ...
//   driveService, err := drive.New(oauthHttpClient)
package drive

import (
	"bytes"
	"code.google.com/p/google-api-go-client/googleapi"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var _ = bytes.NewBuffer
var _ = strconv.Itoa
var _ = fmt.Sprintf
var _ = json.NewDecoder
var _ = io.Copy
var _ = url.Parse
var _ = googleapi.Version
var _ = errors.New

const apiId = "drive:v2"
const apiName = "drive"
const apiVersion = "v2"
const basePath = "https://www.googleapis.com/drive/v2/"

// OAuth2 scopes used by this API.
const (
	// View and manage the files and documents in your Google Drive
	DriveScope = "https://www.googleapis.com/auth/drive"

	// View your Google Drive apps
	DriveAppsReadonlyScope = "https://www.googleapis.com/auth/drive.apps.readonly"

	// View and manage Google Drive files that you have opened or created
	// with this app
	DriveFileScope = "https://www.googleapis.com/auth/drive.file"

	// View metadata for files and documents in your Google Drive
	DriveMetadataReadonlyScope = "https://www.googleapis.com/auth/drive.metadata.readonly"

	// View the files and documents in your Google Drive
	DriveReadonlyScope = "https://www.googleapis.com/auth/drive.readonly"
)

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client}
	s.About = &AboutService{s: s}
	s.Apps = &AppsService{s: s}
	s.Changes = &ChangesService{s: s}
	s.Children = &ChildrenService{s: s}
	s.Comments = &CommentsService{s: s}
	s.Files = &FilesService{s: s}
	s.Parents = &ParentsService{s: s}
	s.Permissions = &PermissionsService{s: s}
	s.Replies = &RepliesService{s: s}
	s.Revisions = &RevisionsService{s: s}
	return s, nil
}

type Service struct {
	client *http.Client

	About *AboutService

	Apps *AppsService

	Changes *ChangesService

	Children *ChildrenService

	Comments *CommentsService

	Files *FilesService

	Parents *ParentsService

	Permissions *PermissionsService

	Replies *RepliesService

	Revisions *RevisionsService
}

type AboutService struct {
	s *Service
}

type AppsService struct {
	s *Service
}

type ChangesService struct {
	s *Service
}

type ChildrenService struct {
	s *Service
}

type CommentsService struct {
	s *Service
}

type FilesService struct {
	s *Service
}

type ParentsService struct {
	s *Service
}

type PermissionsService struct {
	s *Service
}

type RepliesService struct {
	s *Service
}

type RevisionsService struct {
	s *Service
}

type About struct {
	// AdditionalRoleInfo: Information about supported additional roles per
	// file type. The most specific type takes precedence.
	AdditionalRoleInfo []*AboutAdditionalRoleInfo `json:"additionalRoleInfo,omitempty"`

	// DomainSharingPolicy: The domain sharing policy for the current user.
	DomainSharingPolicy string `json:"domainSharingPolicy,omitempty"`

	// Etag: The ETag of the item.
	Etag string `json:"etag,omitempty"`

	// ExportFormats: The allowable export formats.
	ExportFormats []*AboutExportFormats `json:"exportFormats,omitempty"`

	// Features: List of additional features enabled on this account.
	Features []*AboutFeatures `json:"features,omitempty"`

	// ImportFormats: The allowable import formats.
	ImportFormats []*AboutImportFormats `json:"importFormats,omitempty"`

	// IsCurrentAppInstalled: A boolean indicating whether the authenticated
	// app is installed by the authenticated user.
	IsCurrentAppInstalled bool `json:"isCurrentAppInstalled,omitempty"`

	// Kind: This is always drive#about.
	Kind string `json:"kind,omitempty"`

	// LargestChangeId: The largest change id.
	LargestChangeId int64 `json:"largestChangeId,omitempty,string"`

	// MaxUploadSizes: List of max upload sizes for each file type. The most
	// specific type takes precedence.
	MaxUploadSizes []*AboutMaxUploadSizes `json:"maxUploadSizes,omitempty"`

	// Name: The name of the current user.
	Name string `json:"name,omitempty"`

	// PermissionId: The current user's ID as visible in the permissions
	// collection.
	PermissionId string `json:"permissionId,omitempty"`

	// QuotaBytesTotal: The total number of quota bytes.
	QuotaBytesTotal int64 `json:"quotaBytesTotal,omitempty,string"`

	// QuotaBytesUsed: The number of quota bytes used by Google Drive.
	QuotaBytesUsed int64 `json:"quotaBytesUsed,omitempty,string"`

	// QuotaBytesUsedAggregate: The number of quota bytes used by all Google
	// apps (Drive, Picasa, etc.).
	QuotaBytesUsedAggregate int64 `json:"quotaBytesUsedAggregate,omitempty,string"`

	// QuotaBytesUsedInTrash: The number of quota bytes used by trashed
	// items.
	QuotaBytesUsedInTrash int64 `json:"quotaBytesUsedInTrash,omitempty,string"`

	// RemainingChangeIds: The number of remaining change ids.
	RemainingChangeIds int64 `json:"remainingChangeIds,omitempty,string"`

	// RootFolderId: The id of the root folder.
	RootFolderId string `json:"rootFolderId,omitempty"`

	// SelfLink: A link back to this item.
	SelfLink string `json:"selfLink,omitempty"`

	// User: The authenticated user.
	User *User `json:"user,omitempty"`
}

type AboutAdditionalRoleInfo struct {
	// RoleSets: The supported additional roles per primary role.
	RoleSets []*AboutAdditionalRoleInfoRoleSets `json:"roleSets,omitempty"`

	// Type: The content type that this additional role info applies to.
	Type string `json:"type,omitempty"`
}

type AboutAdditionalRoleInfoRoleSets struct {
	// AdditionalRoles: The supported additional roles with the primary
	// role.
	AdditionalRoles []string `json:"additionalRoles,omitempty"`

	// PrimaryRole: A primary permission role.
	PrimaryRole string `json:"primaryRole,omitempty"`
}

type AboutExportFormats struct {
	// Source: The content type to convert from.
	Source string `json:"source,omitempty"`

	// Targets: The possible content types to convert to.
	Targets []string `json:"targets,omitempty"`
}

type AboutFeatures struct {
	// FeatureName: The name of the feature.
	FeatureName string `json:"featureName,omitempty"`

	// FeatureRate: The request limit rate for this feature, in queries per
	// second.
	FeatureRate float64 `json:"featureRate,omitempty"`
}

type AboutImportFormats struct {
	// Source: The imported file's content type to convert from.
	Source string `json:"source,omitempty"`

	// Targets: The possible content types to convert to.
	Targets []string `json:"targets,omitempty"`
}

type AboutMaxUploadSizes struct {
	// Size: The max upload size for this type.
	Size int64 `json:"size,omitempty,string"`

	// Type: The file type.
	Type string `json:"type,omitempty"`
}

type App struct {
	// Authorized: Whether the app is authorized to access data on the
	// user's Drive.
	Authorized bool `json:"authorized,omitempty"`

	// Icons: The various icons for the app.
	Icons []*AppIcons `json:"icons,omitempty"`

	// Id: The ID of the app.
	Id string `json:"id,omitempty"`

	// Installed: Whether the app is installed.
	Installed bool `json:"installed,omitempty"`

	// Kind: This is always drive#app.
	Kind string `json:"kind,omitempty"`

	// Name: The name of the app.
	Name string `json:"name,omitempty"`

	// ObjectType: The type of object this app creates (e.g. Chart). If
	// empty, the app name should be used instead.
	ObjectType string `json:"objectType,omitempty"`

	// PrimaryFileExtensions: The list of primary file extensions.
	PrimaryFileExtensions []string `json:"primaryFileExtensions,omitempty"`

	// PrimaryMimeTypes: The list of primary mime types.
	PrimaryMimeTypes []string `json:"primaryMimeTypes,omitempty"`

	// ProductUrl: The product URL.
	ProductUrl string `json:"productUrl,omitempty"`

	// SecondaryFileExtensions: The list of secondary file extensions.
	SecondaryFileExtensions []string `json:"secondaryFileExtensions,omitempty"`

	// SecondaryMimeTypes: The list of secondary mime types.
	SecondaryMimeTypes []string `json:"secondaryMimeTypes,omitempty"`

	// SupportsCreate: Whether this app supports creating new objects.
	SupportsCreate bool `json:"supportsCreate,omitempty"`

	// SupportsImport: Whether this app supports importing Google Docs.
	SupportsImport bool `json:"supportsImport,omitempty"`

	// UseByDefault: Whether the app is selected as the default handler for
	// the types it supports.
	UseByDefault bool `json:"useByDefault,omitempty"`
}

type AppIcons struct {
	// Category: Category of the icon. Allowed values are:
	// - application -
	// icon for the application
	// - document - icon for a file associated
	// with the app
	// - documentShared - icon for a shared file associated
	// with the app
	Category string `json:"category,omitempty"`

	// IconUrl: URL for the icon.
	IconUrl string `json:"iconUrl,omitempty"`

	// Size: Size of the icon. Represented as the maximum of the width and
	// height.
	Size int64 `json:"size,omitempty"`
}

type AppList struct {
	// Etag: The ETag of the list.
	Etag string `json:"etag,omitempty"`

	// Items: The actual list of apps.
	Items []*App `json:"items,omitempty"`

	// Kind: This is always drive#appList.
	Kind string `json:"kind,omitempty"`

	// SelfLink: A link back to this list.
	SelfLink string `json:"selfLink,omitempty"`
}

type Change struct {
	// Deleted: Whether the file has been deleted.
	Deleted bool `json:"deleted,omitempty"`

	// File: The updated state of the file. Present if the file has not been
	// deleted.
	File *File `json:"file,omitempty"`

	// FileId: The ID of the file associated with this change.
	FileId string `json:"fileId,omitempty"`

	// Id: The ID of the change.
	Id int64 `json:"id,omitempty,string"`

	// Kind: This is always drive#change.
	Kind string `json:"kind,omitempty"`

	// SelfLink: A link back to this change.
	SelfLink string `json:"selfLink,omitempty"`
}

type ChangeList struct {
	// Etag: The ETag of the list.
	Etag string `json:"etag,omitempty"`

	// Items: The actual list of changes.
	Items []*Change `json:"items,omitempty"`

	// Kind: This is always drive#changeList.
	Kind string `json:"kind,omitempty"`

	// LargestChangeId: The current largest change ID.
	LargestChangeId int64 `json:"largestChangeId,omitempty,string"`

	// NextLink: A link to the next page of changes.
	NextLink string `json:"nextLink,omitempty"`

	// NextPageToken: The page token for the next page of changes.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// SelfLink: A link back to this list.
	SelfLink string `json:"selfLink,omitempty"`
}

type ChildList struct {
	// Etag: The ETag of the list.
	Etag string `json:"etag,omitempty"`

	// Items: The actual list of children.
	Items []*ChildReference `json:"items,omitempty"`

	// Kind: This is always drive#childList.
	Kind string `json:"kind,omitempty"`

	// NextLink: A link to the next page of children.
	NextLink string `json:"nextLink,omitempty"`

	// NextPageToken: The page token for the next page of children.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// SelfLink: A link back to this list.
	SelfLink string `json:"selfLink,omitempty"`
}

type ChildReference struct {
	// ChildLink: A link to the child.
	ChildLink string `json:"childLink,omitempty"`

	// Id: The ID of the child.
	Id string `json:"id,omitempty"`

	// Kind: This is always drive#childReference.
	Kind string `json:"kind,omitempty"`

	// SelfLink: A link back to this reference.
	SelfLink string `json:"selfLink,omitempty"`
}

type Comment struct {
	// Anchor: A region of the document represented as a JSON string. See
	// anchor documentation for details on how to define and interpret
	// anchor properties.
	Anchor string `json:"anchor,omitempty"`

	// Author: The user who wrote this comment.
	Author *User `json:"author,omitempty"`

	// CommentId: The ID of the comment.
	CommentId string `json:"commentId,omitempty"`

	// Content: The plain text content used to create this comment. This is
	// not HTML safe and should only be used as a starting point to make
	// edits to a comment's content.
	Content string `json:"content,omitempty"`

	// Context: The context of the file which is being commented on.
	Context *CommentContext `json:"context,omitempty"`

	// CreatedDate: The date when this comment was first created.
	CreatedDate string `json:"createdDate,omitempty"`

	// Deleted: Whether this comment has been deleted. If a comment has been
	// deleted the content will be cleared and this will only represent a
	// comment that once existed.
	Deleted bool `json:"deleted,omitempty"`

	// FileId: The file which this comment is addressing.
	FileId string `json:"fileId,omitempty"`

	// FileTitle: The title of the file which this comment is addressing.
	FileTitle string `json:"fileTitle,omitempty"`

	// HtmlContent: HTML formatted content for this comment.
	HtmlContent string `json:"htmlContent,omitempty"`

	// Kind: This is always drive#comment.
	Kind string `json:"kind,omitempty"`

	// ModifiedDate: The date when this comment or any of its replies were
	// last modified.
	ModifiedDate string `json:"modifiedDate,omitempty"`

	// Replies: Replies to this post.
	Replies []*CommentReply `json:"replies,omitempty"`

	// SelfLink: A link back to this comment.
	SelfLink string `json:"selfLink,omitempty"`

	// Status: The status of this comment. Status can be changed by posting
	// a reply to a comment with the desired status.
	// - "open" - The
	// comment is still open.
	// - "resolved" - The comment has been resolved
	// by one of its replies.
	Status string `json:"status,omitempty"`
}

type CommentContext struct {
	// Type: The MIME type of the context snippet.
	Type string `json:"type,omitempty"`

	// Value: Data representation of the segment of the file being commented
	// on. In the case of a text file for example, this would be the actual
	// text that the comment is about.
	Value string `json:"value,omitempty"`
}

type CommentList struct {
	// Items: List of comments.
	Items []*Comment `json:"items,omitempty"`

	// Kind: This is always drive#commentList.
	Kind string `json:"kind,omitempty"`

	// NextPageToken: The token to use to request the next page of results.
	NextPageToken string `json:"nextPageToken,omitempty"`
}

type CommentReply struct {
	// Author: The user who wrote this reply.
	Author *User `json:"author,omitempty"`

	// Content: The plain text content used to create this reply. This is
	// not HTML safe and should only be used as a starting point to make
	// edits to a reply's content. This field is required on inserts if no
	// verb is specified (resolve/reopen).
	Content string `json:"content,omitempty"`

	// CreatedDate: The date when this reply was first created.
	CreatedDate string `json:"createdDate,omitempty"`

	// Deleted: Whether this reply has been deleted. If a reply has been
	// deleted the content will be cleared and this will only represent a
	// reply that once existed.
	Deleted bool `json:"deleted,omitempty"`

	// HtmlContent: HTML formatted content for this reply.
	HtmlContent string `json:"htmlContent,omitempty"`

	// Kind: This is always drive#commentReply.
	Kind string `json:"kind,omitempty"`

	// ModifiedDate: The date when this reply was last modified.
	ModifiedDate string `json:"modifiedDate,omitempty"`

	// ReplyId: The ID of the reply.
	ReplyId string `json:"replyId,omitempty"`

	// Verb: The action this reply performed to the parent comment. When
	// creating a new reply this is the action to be perform to the parent
	// comment. Possible values are:
	// - "resolve" - To resolve a comment.
	//
	// - "reopen" - To reopen (un-resolve) a comment.
	Verb string `json:"verb,omitempty"`
}

type CommentReplyList struct {
	// Items: List of reply.
	Items []*CommentReply `json:"items,omitempty"`

	// Kind: This is always drive#commentReplyList.
	Kind string `json:"kind,omitempty"`

	// NextPageToken: The token to use to request the next page of results.
	NextPageToken string `json:"nextPageToken,omitempty"`
}

type File struct {
	// AlternateLink: A link for opening the file in using a relevant Google
	// editor or viewer.
	AlternateLink string `json:"alternateLink,omitempty"`

	// CreatedDate: Create time for this file (formatted ISO8601 timestamp).
	CreatedDate string `json:"createdDate,omitempty"`

	// Description: A short description of the file.
	Description string `json:"description,omitempty"`

	// DownloadUrl: Short lived download URL for the file. This is only
	// populated for files with content stored in Drive.
	DownloadUrl string `json:"downloadUrl,omitempty"`

	// Editable: Whether the file can be edited by the current user.
	Editable bool `json:"editable,omitempty"`

	// EmbedLink: A link for embedding the file.
	EmbedLink string `json:"embedLink,omitempty"`

	// Etag: ETag of the file.
	Etag string `json:"etag,omitempty"`

	// ExplicitlyTrashed: Whether this file has been explicitly trashed, as
	// opposed to recursively trashed. This will only be populated if the
	// file is trashed.
	ExplicitlyTrashed bool `json:"explicitlyTrashed,omitempty"`

	// ExportLinks: Links for exporting Google Docs to specific formats.
	ExportLinks *FileExportLinks `json:"exportLinks,omitempty"`

	// FileExtension: The file extension used when downloading this file.
	// This field is read only. To set the extension, include it in the
	// title when creating the file. This is only populated for files with
	// content stored in Drive.
	FileExtension string `json:"fileExtension,omitempty"`

	// FileSize: The size of the file in bytes. This is only populated for
	// files with content stored in Drive.
	FileSize int64 `json:"fileSize,omitempty,string"`

	// IconLink: A link to the file's icon.
	IconLink string `json:"iconLink,omitempty"`

	// Id: The ID of the file.
	Id string `json:"id,omitempty"`

	// ImageMediaMetadata: Metadata about image media. This will only be
	// present for image types, and its contents will depend on what can be
	// parsed from the image content.
	ImageMediaMetadata *FileImageMediaMetadata `json:"imageMediaMetadata,omitempty"`

	// IndexableText: Indexable text attributes for the file (can only be
	// written)
	IndexableText *FileIndexableText `json:"indexableText,omitempty"`

	// Kind: The type of file. This is always drive#file.
	Kind string `json:"kind,omitempty"`

	// Labels: A group of labels for the file.
	Labels *FileLabels `json:"labels,omitempty"`

	// LastModifyingUser: The last user to modify this file.
	LastModifyingUser *User `json:"lastModifyingUser,omitempty"`

	// LastModifyingUserName: Name of the last user to modify this file.
	LastModifyingUserName string `json:"lastModifyingUserName,omitempty"`

	// LastViewedByMeDate: Last time this file was viewed by the user
	// (formatted RFC 3339 timestamp).
	LastViewedByMeDate string `json:"lastViewedByMeDate,omitempty"`

	// Md5Checksum: An MD5 checksum for the content of this file. This is
	// populated only for files with content stored in Drive.
	Md5Checksum string `json:"md5Checksum,omitempty"`

	// MimeType: The MIME type of the file. This is only mutable on update
	// when uploading new content. This field can be left blank, and the
	// mimetype will be determined from the uploaded content's MIME type.
	MimeType string `json:"mimeType,omitempty"`

	// ModifiedByMeDate: Last time this file was modified by the user
	// (formatted RFC 3339 timestamp). Note that setting modifiedDate will
	// also update the modifiedByMe date for the user which set the date.
	ModifiedByMeDate string `json:"modifiedByMeDate,omitempty"`

	// ModifiedDate: Last time this file was modified by anyone (formatted
	// RFC 3339 timestamp). This is only mutable on update when the
	// setModifiedDate parameter is set.
	ModifiedDate string `json:"modifiedDate,omitempty"`

	// OriginalFilename: The original filename if the file was uploaded
	// manually, or the original title if the file was inserted through the
	// API. Note that renames of the title will not change the original
	// filename. This will only be populated on files with content stored in
	// Drive.
	OriginalFilename string `json:"originalFilename,omitempty"`

	// OwnerNames: Name(s) of the owner(s) of this file.
	OwnerNames []string `json:"ownerNames,omitempty"`

	// Owners: The owner(s) of this file.
	Owners []*User `json:"owners,omitempty"`

	// Parents: Collection of parent folders which contain this
	// file.
	// Setting this field will put the file in all of the provided
	// folders. On insert, if no folders are provided, the file will be
	// placed in the default root folder.
	Parents []*ParentReference `json:"parents,omitempty"`

	// QuotaBytesUsed: The number of quota bytes used by this file.
	QuotaBytesUsed int64 `json:"quotaBytesUsed,omitempty,string"`

	// SelfLink: A link back to this file.
	SelfLink string `json:"selfLink,omitempty"`

	// Shared: Whether the file has been shared.
	Shared bool `json:"shared,omitempty"`

	// SharedWithMeDate: Time at which this file was shared with the user
	// (formatted RFC 3339 timestamp).
	SharedWithMeDate string `json:"sharedWithMeDate,omitempty"`

	// Thumbnail: Thumbnail for the file. Only accepted on upload and for
	// files that are not already thumbnailed by Google.
	Thumbnail *FileThumbnail `json:"thumbnail,omitempty"`

	// ThumbnailLink: A link to the file's thumbnail.
	ThumbnailLink string `json:"thumbnailLink,omitempty"`

	// Title: The title of this file.
	Title string `json:"title,omitempty"`

	// UserPermission: The permissions for the authenticated user on this
	// file.
	UserPermission *Permission `json:"userPermission,omitempty"`

	// WebContentLink: A link for downloading the content of the file in a
	// browser using cookie based authentication. In cases where the content
	// is shared publicly, the content can be downloaded without any
	// credentials.
	WebContentLink string `json:"webContentLink,omitempty"`

	// WebViewLink: A link only available on public folders for viewing
	// their static web assets (HTML, CSS, JS, etc) via Google Drive's
	// Website Hosting.
	WebViewLink string `json:"webViewLink,omitempty"`

	// WritersCanShare: Whether writers can share the document with other
	// users.
	WritersCanShare bool `json:"writersCanShare,omitempty"`
}

type FileExportLinks struct {
}

type FileImageMediaMetadata struct {
	// Aperture: The aperture used to create the photo (f-number).
	Aperture float64 `json:"aperture,omitempty"`

	// CameraMake: The make of the camera used to create the photo.
	CameraMake string `json:"cameraMake,omitempty"`

	// CameraModel: The model of the camera used to create the photo.
	CameraModel string `json:"cameraModel,omitempty"`

	// ColorSpace: The color space of the photo.
	ColorSpace string `json:"colorSpace,omitempty"`

	// Date: The date and time the photo was taken (EXIF format timestamp).
	Date string `json:"date,omitempty"`

	// ExposureBias: The exposure bias of the photo (APEX value).
	ExposureBias float64 `json:"exposureBias,omitempty"`

	// ExposureMode: The exposure mode used to create the photo.
	ExposureMode string `json:"exposureMode,omitempty"`

	// ExposureTime: The length of the exposure, in seconds.
	ExposureTime float64 `json:"exposureTime,omitempty"`

	// FlashUsed: Whether a flash was used to create the photo.
	FlashUsed bool `json:"flashUsed,omitempty"`

	// FocalLength: The focal length used to create the photo, in
	// millimeters.
	FocalLength float64 `json:"focalLength,omitempty"`

	// Height: The height of the image in pixels.
	Height int64 `json:"height,omitempty"`

	// IsoSpeed: The ISO speed used to create the photo.
	IsoSpeed int64 `json:"isoSpeed,omitempty"`

	// Lens: The lens used to create the photo.
	Lens string `json:"lens,omitempty"`

	// Location: Geographic location information stored in the image.
	Location *FileImageMediaMetadataLocation `json:"location,omitempty"`

	// MaxApertureValue: The smallest f-number of the lens at the focal
	// length used to create the photo (APEX value).
	MaxApertureValue float64 `json:"maxApertureValue,omitempty"`

	// MeteringMode: The metering mode used to create the photo.
	MeteringMode string `json:"meteringMode,omitempty"`

	// Rotation: The rotation in clockwise degrees from the image's original
	// orientation.
	Rotation int64 `json:"rotation,omitempty"`

	// Sensor: The type of sensor used to create the photo.
	Sensor string `json:"sensor,omitempty"`

	// SubjectDistance: The distance to the subject of the photo, in meters.
	SubjectDistance int64 `json:"subjectDistance,omitempty"`

	// WhiteBalance: The white balance mode used to create the photo.
	WhiteBalance string `json:"whiteBalance,omitempty"`

	// Width: The width of the image in pixels.
	Width int64 `json:"width,omitempty"`
}

type FileImageMediaMetadataLocation struct {
	// Altitude: The altitude stored in the image.
	Altitude float64 `json:"altitude,omitempty"`

	// Latitude: The latitude stored in the image.
	Latitude float64 `json:"latitude,omitempty"`

	// Longitude: The longitude stored in the image.
	Longitude float64 `json:"longitude,omitempty"`
}

type FileIndexableText struct {
	// Text: The text to be indexed for this file.
	Text string `json:"text,omitempty"`
}

type FileLabels struct {
	// Hidden: Whether this file is hidden from the user.
	Hidden bool `json:"hidden,omitempty"`

	// Restricted: Whether viewers are prevented from downloading this file.
	Restricted bool `json:"restricted,omitempty"`

	// Starred: Whether this file is starred by the user.
	Starred bool `json:"starred,omitempty"`

	// Trashed: Whether this file has been trashed.
	Trashed bool `json:"trashed,omitempty"`

	// Viewed: Whether this file has been viewed by this user.
	Viewed bool `json:"viewed,omitempty"`
}

type FileThumbnail struct {
	// Image: The URL-safe Base64 encoded bytes of the thumbnail image.
	Image string `json:"image,omitempty"`

	// MimeType: The MIME type of the thumbnail.
	MimeType string `json:"mimeType,omitempty"`
}

type FileList struct {
	// Etag: The ETag of the list.
	Etag string `json:"etag,omitempty"`

	// Items: The actual list of files.
	Items []*File `json:"items,omitempty"`

	// Kind: This is always drive#fileList.
	Kind string `json:"kind,omitempty"`

	// NextLink: A link to the next page of files.
	NextLink string `json:"nextLink,omitempty"`

	// NextPageToken: The page token for the next page of files.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// SelfLink: A link back to this list.
	SelfLink string `json:"selfLink,omitempty"`
}

type ParentList struct {
	// Etag: The ETag of the list.
	Etag string `json:"etag,omitempty"`

	// Items: The actual list of parents.
	Items []*ParentReference `json:"items,omitempty"`

	// Kind: This is always drive#parentList.
	Kind string `json:"kind,omitempty"`

	// SelfLink: A link back to this list.
	SelfLink string `json:"selfLink,omitempty"`
}

type ParentReference struct {
	// Id: The ID of the parent.
	Id string `json:"id,omitempty"`

	// IsRoot: Whether or not the parent is the root folder.
	IsRoot bool `json:"isRoot,omitempty"`

	// Kind: This is always drive#parentReference.
	Kind string `json:"kind,omitempty"`

	// ParentLink: A link to the parent.
	ParentLink string `json:"parentLink,omitempty"`

	// SelfLink: A link back to this reference.
	SelfLink string `json:"selfLink,omitempty"`
}

type Permission struct {
	// AdditionalRoles: Additional roles for this user. Only commenter is
	// currently allowed.
	AdditionalRoles []string `json:"additionalRoles,omitempty"`

	// AuthKey: The authkey parameter required for this permission.
	AuthKey string `json:"authKey,omitempty"`

	// Etag: The ETag of the permission.
	Etag string `json:"etag,omitempty"`

	// Id: The ID of the permission.
	Id string `json:"id,omitempty"`

	// Kind: This is always drive#permission.
	Kind string `json:"kind,omitempty"`

	// Name: The name for this permission.
	Name string `json:"name,omitempty"`

	// PhotoLink: A link to the profile photo, if available.
	PhotoLink string `json:"photoLink,omitempty"`

	// Role: The primary role for this user. Allowed values are:
	// - owner
	//
	// - reader
	// - writer
	Role string `json:"role,omitempty"`

	// SelfLink: A link back to this permission.
	SelfLink string `json:"selfLink,omitempty"`

	// Type: The account type. Allowed values are:
	// - user
	// - group
	// -
	// domain
	// - anyone
	Type string `json:"type,omitempty"`

	// Value: The email address or domain name for the entity. This is not
	// populated in responses.
	Value string `json:"value,omitempty"`

	// WithLink: Whether the link is required for this permission.
	WithLink bool `json:"withLink,omitempty"`
}

type PermissionList struct {
	// Etag: The ETag of the list.
	Etag string `json:"etag,omitempty"`

	// Items: The actual list of permissions.
	Items []*Permission `json:"items,omitempty"`

	// Kind: This is always drive#permissionList.
	Kind string `json:"kind,omitempty"`

	// SelfLink: A link back to this list.
	SelfLink string `json:"selfLink,omitempty"`
}

type Revision struct {
	// DownloadUrl: Short term download URL for the file. This will only be
	// populated on files with content stored in Drive.
	DownloadUrl string `json:"downloadUrl,omitempty"`

	// Etag: The ETag of the revision.
	Etag string `json:"etag,omitempty"`

	// ExportLinks: Links for exporting Google Docs to specific formats.
	ExportLinks *RevisionExportLinks `json:"exportLinks,omitempty"`

	// FileSize: The size of the revision in bytes. This will only be
	// populated on files with content stored in Drive.
	FileSize int64 `json:"fileSize,omitempty,string"`

	// Id: The ID of the revision.
	Id string `json:"id,omitempty"`

	// Kind: This is always drive#revision.
	Kind string `json:"kind,omitempty"`

	// LastModifyingUser: The last user to modify this revision.
	LastModifyingUser *User `json:"lastModifyingUser,omitempty"`

	// LastModifyingUserName: Name of the last user to modify this revision.
	LastModifyingUserName string `json:"lastModifyingUserName,omitempty"`

	// Md5Checksum: An MD5 checksum for the content of this revision. This
	// will only be populated on files with content stored in Drive.
	Md5Checksum string `json:"md5Checksum,omitempty"`

	// MimeType: The MIME type of the revision.
	MimeType string `json:"mimeType,omitempty"`

	// ModifiedDate: Last time this revision was modified (formatted RFC
	// 3339 timestamp).
	ModifiedDate string `json:"modifiedDate,omitempty"`

	// OriginalFilename: The original filename when this revision was
	// created. This will only be populated on files with content stored in
	// Drive.
	OriginalFilename string `json:"originalFilename,omitempty"`

	// Pinned: Whether this revision is pinned to prevent automatic purging.
	// This will only be populated and can only be modified on files with
	// content stored in Drive which are not Google Docs. Revisions can also
	// be pinned when they are created through the
	// drive.files.insert/update/copy by using the pinned query parameter.
	Pinned bool `json:"pinned,omitempty"`

	// PublishAuto: Whether subsequent revisions will be automatically
	// republished. This is only populated and can only be modified for
	// Google Docs.
	PublishAuto bool `json:"publishAuto,omitempty"`

	// Published: Whether this revision is published. This is only populated
	// and can only be modified for Google Docs.
	Published bool `json:"published,omitempty"`

	// PublishedLink: A link to the published revision.
	PublishedLink string `json:"publishedLink,omitempty"`

	// PublishedOutsideDomain: Whether this revision is published outside
	// the domain. This is only populated and can only be modified for
	// Google Docs.
	PublishedOutsideDomain bool `json:"publishedOutsideDomain,omitempty"`

	// SelfLink: A link back to this revision.
	SelfLink string `json:"selfLink,omitempty"`
}

type RevisionExportLinks struct {
}

type RevisionList struct {
	// Etag: The ETag of the list.
	Etag string `json:"etag,omitempty"`

	// Items: The actual list of revisions.
	Items []*Revision `json:"items,omitempty"`

	// Kind: This is always drive#revisionList.
	Kind string `json:"kind,omitempty"`

	// SelfLink: A link back to this list.
	SelfLink string `json:"selfLink,omitempty"`
}

type User struct {
	// DisplayName: A plain text displayable name for this user.
	DisplayName string `json:"displayName,omitempty"`

	// IsAuthenticatedUser: Whether this user is the same as the
	// authenticated user for whom the request was made.
	IsAuthenticatedUser bool `json:"isAuthenticatedUser,omitempty"`

	// Kind: This is always drive#user.
	Kind string `json:"kind,omitempty"`

	// PermissionId: The user's ID as visible in the permissions collection.
	PermissionId string `json:"permissionId,omitempty"`

	// Picture: The user's profile picture.
	Picture *UserPicture `json:"picture,omitempty"`
}

type UserPicture struct {
	// Url: A URL that points to a profile picture of this user.
	Url string `json:"url,omitempty"`
}

// method id "drive.about.get":

type AboutGetCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// Get: Gets the information about the current user along with Drive API
// settings
func (r *AboutService) Get() *AboutGetCall {
	c := &AboutGetCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

// IncludeSubscribed sets the optional parameter "includeSubscribed":
// When calculating the number of remaining change IDs, whether to
// include shared files and public files the user has opened. When set
// to false, this counts only change IDs for owned files and any shared
// or public files that the user has explictly added to a folder in
// Drive.
func (c *AboutGetCall) IncludeSubscribed(includeSubscribed bool) *AboutGetCall {
	c.opt_["includeSubscribed"] = includeSubscribed
	return c
}

// MaxChangeIdCount sets the optional parameter "maxChangeIdCount":
// Maximum number of remaining change IDs to count
func (c *AboutGetCall) MaxChangeIdCount(maxChangeIdCount int64) *AboutGetCall {
	c.opt_["maxChangeIdCount"] = maxChangeIdCount
	return c
}

// StartChangeId sets the optional parameter "startChangeId": Change ID
// to start counting from when calculating number of remaining change
// IDs
func (c *AboutGetCall) StartChangeId(startChangeId int64) *AboutGetCall {
	c.opt_["startChangeId"] = startChangeId
	return c
}

func (c *AboutGetCall) Do() (*About, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["includeSubscribed"]; ok {
		params.Set("includeSubscribed", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["maxChangeIdCount"]; ok {
		params.Set("maxChangeIdCount", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["startChangeId"]; ok {
		params.Set("startChangeId", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "about")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(About)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets the information about the current user along with Drive API settings",
	//   "httpMethod": "GET",
	//   "id": "drive.about.get",
	//   "parameters": {
	//     "includeSubscribed": {
	//       "default": "true",
	//       "description": "When calculating the number of remaining change IDs, whether to include shared files and public files the user has opened. When set to false, this counts only change IDs for owned files and any shared or public files that the user has explictly added to a folder in Drive.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "maxChangeIdCount": {
	//       "default": "1",
	//       "description": "Maximum number of remaining change IDs to count",
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "startChangeId": {
	//       "description": "Change ID to start counting from when calculating number of remaining change IDs",
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "about",
	//   "response": {
	//     "$ref": "About"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.apps.get":

type AppsGetCall struct {
	s     *Service
	appId string
	opt_  map[string]interface{}
}

// Get: Gets a specific app.
func (r *AppsService) Get(appId string) *AppsGetCall {
	c := &AppsGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.appId = appId
	return c
}

func (c *AppsGetCall) Do() (*App, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "apps/{appId}")
	urls = strings.Replace(urls, "{appId}", cleanPathString(c.appId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(App)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets a specific app.",
	//   "httpMethod": "GET",
	//   "id": "drive.apps.get",
	//   "parameterOrder": [
	//     "appId"
	//   ],
	//   "parameters": {
	//     "appId": {
	//       "description": "The ID of the app.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "apps/{appId}",
	//   "response": {
	//     "$ref": "App"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive.apps.readonly"
	//   ]
	// }

}

// method id "drive.apps.list":

type AppsListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// List: Lists a user's installed apps.
func (r *AppsService) List() *AppsListCall {
	c := &AppsListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

func (c *AppsListCall) Do() (*AppList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "apps")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(AppList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Lists a user's installed apps.",
	//   "httpMethod": "GET",
	//   "id": "drive.apps.list",
	//   "path": "apps",
	//   "response": {
	//     "$ref": "AppList"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive.apps.readonly"
	//   ]
	// }

}

// method id "drive.changes.get":

type ChangesGetCall struct {
	s        *Service
	changeId string
	opt_     map[string]interface{}
}

// Get: Gets a specific change.
func (r *ChangesService) Get(changeId string) *ChangesGetCall {
	c := &ChangesGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.changeId = changeId
	return c
}

func (c *ChangesGetCall) Do() (*Change, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "changes/{changeId}")
	urls = strings.Replace(urls, "{changeId}", cleanPathString(c.changeId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Change)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets a specific change.",
	//   "httpMethod": "GET",
	//   "id": "drive.changes.get",
	//   "parameterOrder": [
	//     "changeId"
	//   ],
	//   "parameters": {
	//     "changeId": {
	//       "description": "The ID of the change.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "changes/{changeId}",
	//   "response": {
	//     "$ref": "Change"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.changes.list":

type ChangesListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// List: Lists the changes for a user.
func (r *ChangesService) List() *ChangesListCall {
	c := &ChangesListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

// IncludeDeleted sets the optional parameter "includeDeleted": Whether
// to include deleted items.
func (c *ChangesListCall) IncludeDeleted(includeDeleted bool) *ChangesListCall {
	c.opt_["includeDeleted"] = includeDeleted
	return c
}

// IncludeSubscribed sets the optional parameter "includeSubscribed":
// Whether to include shared files and public files the user has opened.
// When set to false, the list will include owned files plus any shared
// or public files the user has explictly added to a folder in Drive.
func (c *ChangesListCall) IncludeSubscribed(includeSubscribed bool) *ChangesListCall {
	c.opt_["includeSubscribed"] = includeSubscribed
	return c
}

// MaxResults sets the optional parameter "maxResults": Maximum number
// of changes to return.
func (c *ChangesListCall) MaxResults(maxResults int64) *ChangesListCall {
	c.opt_["maxResults"] = maxResults
	return c
}

// PageToken sets the optional parameter "pageToken": Page token for
// changes.
func (c *ChangesListCall) PageToken(pageToken string) *ChangesListCall {
	c.opt_["pageToken"] = pageToken
	return c
}

// StartChangeId sets the optional parameter "startChangeId": Change ID
// to start listing changes from.
func (c *ChangesListCall) StartChangeId(startChangeId int64) *ChangesListCall {
	c.opt_["startChangeId"] = startChangeId
	return c
}

func (c *ChangesListCall) Do() (*ChangeList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["includeDeleted"]; ok {
		params.Set("includeDeleted", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["includeSubscribed"]; ok {
		params.Set("includeSubscribed", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["maxResults"]; ok {
		params.Set("maxResults", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["pageToken"]; ok {
		params.Set("pageToken", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["startChangeId"]; ok {
		params.Set("startChangeId", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "changes")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(ChangeList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Lists the changes for a user.",
	//   "httpMethod": "GET",
	//   "id": "drive.changes.list",
	//   "parameters": {
	//     "includeDeleted": {
	//       "default": "true",
	//       "description": "Whether to include deleted items.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "includeSubscribed": {
	//       "default": "true",
	//       "description": "Whether to include shared files and public files the user has opened. When set to false, the list will include owned files plus any shared or public files the user has explictly added to a folder in Drive.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "maxResults": {
	//       "default": "100",
	//       "description": "Maximum number of changes to return.",
	//       "format": "int32",
	//       "location": "query",
	//       "minimum": "0",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Page token for changes.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "startChangeId": {
	//       "description": "Change ID to start listing changes from.",
	//       "format": "int64",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "changes",
	//   "response": {
	//     "$ref": "ChangeList"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ],
	//   "supportsSubscription": true
	// }

}

// method id "drive.children.delete":

type ChildrenDeleteCall struct {
	s        *Service
	folderId string
	childId  string
	opt_     map[string]interface{}
}

// Delete: Removes a child from a folder.
func (r *ChildrenService) Delete(folderId string, childId string) *ChildrenDeleteCall {
	c := &ChildrenDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.folderId = folderId
	c.childId = childId
	return c
}

func (c *ChildrenDeleteCall) Do() error {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{folderId}/children/{childId}")
	urls = strings.Replace(urls, "{folderId}", cleanPathString(c.folderId), 1)
	urls = strings.Replace(urls, "{childId}", cleanPathString(c.childId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
	// {
	//   "description": "Removes a child from a folder.",
	//   "httpMethod": "DELETE",
	//   "id": "drive.children.delete",
	//   "parameterOrder": [
	//     "folderId",
	//     "childId"
	//   ],
	//   "parameters": {
	//     "childId": {
	//       "description": "The ID of the child.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "folderId": {
	//       "description": "The ID of the folder.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{folderId}/children/{childId}",
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.children.get":

type ChildrenGetCall struct {
	s        *Service
	folderId string
	childId  string
	opt_     map[string]interface{}
}

// Get: Gets a specific child reference.
func (r *ChildrenService) Get(folderId string, childId string) *ChildrenGetCall {
	c := &ChildrenGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.folderId = folderId
	c.childId = childId
	return c
}

func (c *ChildrenGetCall) Do() (*ChildReference, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{folderId}/children/{childId}")
	urls = strings.Replace(urls, "{folderId}", cleanPathString(c.folderId), 1)
	urls = strings.Replace(urls, "{childId}", cleanPathString(c.childId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(ChildReference)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets a specific child reference.",
	//   "httpMethod": "GET",
	//   "id": "drive.children.get",
	//   "parameterOrder": [
	//     "folderId",
	//     "childId"
	//   ],
	//   "parameters": {
	//     "childId": {
	//       "description": "The ID of the child.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "folderId": {
	//       "description": "The ID of the folder.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{folderId}/children/{childId}",
	//   "response": {
	//     "$ref": "ChildReference"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.children.insert":

type ChildrenInsertCall struct {
	s              *Service
	folderId       string
	childreference *ChildReference
	opt_           map[string]interface{}
}

// Insert: Inserts a file into a folder.
func (r *ChildrenService) Insert(folderId string, childreference *ChildReference) *ChildrenInsertCall {
	c := &ChildrenInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.folderId = folderId
	c.childreference = childreference
	return c
}

func (c *ChildrenInsertCall) Do() (*ChildReference, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.childreference)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{folderId}/children")
	urls = strings.Replace(urls, "{folderId}", cleanPathString(c.folderId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(ChildReference)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Inserts a file into a folder.",
	//   "httpMethod": "POST",
	//   "id": "drive.children.insert",
	//   "parameterOrder": [
	//     "folderId"
	//   ],
	//   "parameters": {
	//     "folderId": {
	//       "description": "The ID of the folder.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{folderId}/children",
	//   "request": {
	//     "$ref": "ChildReference"
	//   },
	//   "response": {
	//     "$ref": "ChildReference"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.children.list":

type ChildrenListCall struct {
	s        *Service
	folderId string
	opt_     map[string]interface{}
}

// List: Lists a folder's children.
func (r *ChildrenService) List(folderId string) *ChildrenListCall {
	c := &ChildrenListCall{s: r.s, opt_: make(map[string]interface{})}
	c.folderId = folderId
	return c
}

// MaxResults sets the optional parameter "maxResults": Maximum number
// of children to return.
func (c *ChildrenListCall) MaxResults(maxResults int64) *ChildrenListCall {
	c.opt_["maxResults"] = maxResults
	return c
}

// PageToken sets the optional parameter "pageToken": Page token for
// children.
func (c *ChildrenListCall) PageToken(pageToken string) *ChildrenListCall {
	c.opt_["pageToken"] = pageToken
	return c
}

// Q sets the optional parameter "q": Query string for searching
// children.
func (c *ChildrenListCall) Q(q string) *ChildrenListCall {
	c.opt_["q"] = q
	return c
}

func (c *ChildrenListCall) Do() (*ChildList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["maxResults"]; ok {
		params.Set("maxResults", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["pageToken"]; ok {
		params.Set("pageToken", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["q"]; ok {
		params.Set("q", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{folderId}/children")
	urls = strings.Replace(urls, "{folderId}", cleanPathString(c.folderId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(ChildList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Lists a folder's children.",
	//   "httpMethod": "GET",
	//   "id": "drive.children.list",
	//   "parameterOrder": [
	//     "folderId"
	//   ],
	//   "parameters": {
	//     "folderId": {
	//       "description": "The ID of the folder.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "maxResults": {
	//       "default": "100",
	//       "description": "Maximum number of children to return.",
	//       "format": "int32",
	//       "location": "query",
	//       "minimum": "0",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Page token for children.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "q": {
	//       "description": "Query string for searching children.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{folderId}/children",
	//   "response": {
	//     "$ref": "ChildList"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.comments.delete":

type CommentsDeleteCall struct {
	s         *Service
	fileId    string
	commentId string
	opt_      map[string]interface{}
}

// Delete: Deletes a comment.
func (r *CommentsService) Delete(fileId string, commentId string) *CommentsDeleteCall {
	c := &CommentsDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.commentId = commentId
	return c
}

func (c *CommentsDeleteCall) Do() error {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments/{commentId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{commentId}", cleanPathString(c.commentId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
	// {
	//   "description": "Deletes a comment.",
	//   "httpMethod": "DELETE",
	//   "id": "drive.comments.delete",
	//   "parameterOrder": [
	//     "fileId",
	//     "commentId"
	//   ],
	//   "parameters": {
	//     "commentId": {
	//       "description": "The ID of the comment.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments/{commentId}",
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.comments.get":

type CommentsGetCall struct {
	s         *Service
	fileId    string
	commentId string
	opt_      map[string]interface{}
}

// Get: Gets a comment by ID.
func (r *CommentsService) Get(fileId string, commentId string) *CommentsGetCall {
	c := &CommentsGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.commentId = commentId
	return c
}

// IncludeDeleted sets the optional parameter "includeDeleted": If set,
// this will succeed when retrieving a deleted comment, and will include
// any deleted replies.
func (c *CommentsGetCall) IncludeDeleted(includeDeleted bool) *CommentsGetCall {
	c.opt_["includeDeleted"] = includeDeleted
	return c
}

func (c *CommentsGetCall) Do() (*Comment, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["includeDeleted"]; ok {
		params.Set("includeDeleted", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments/{commentId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{commentId}", cleanPathString(c.commentId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Comment)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets a comment by ID.",
	//   "httpMethod": "GET",
	//   "id": "drive.comments.get",
	//   "parameterOrder": [
	//     "fileId",
	//     "commentId"
	//   ],
	//   "parameters": {
	//     "commentId": {
	//       "description": "The ID of the comment.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "includeDeleted": {
	//       "default": "false",
	//       "description": "If set, this will succeed when retrieving a deleted comment, and will include any deleted replies.",
	//       "location": "query",
	//       "type": "boolean"
	//     }
	//   },
	//   "path": "files/{fileId}/comments/{commentId}",
	//   "response": {
	//     "$ref": "Comment"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.comments.insert":

type CommentsInsertCall struct {
	s       *Service
	fileId  string
	comment *Comment
	opt_    map[string]interface{}
}

// Insert: Creates a new comment on the given file.
func (r *CommentsService) Insert(fileId string, comment *Comment) *CommentsInsertCall {
	c := &CommentsInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.comment = comment
	return c
}

func (c *CommentsInsertCall) Do() (*Comment, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.comment)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Comment)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates a new comment on the given file.",
	//   "httpMethod": "POST",
	//   "id": "drive.comments.insert",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments",
	//   "request": {
	//     "$ref": "Comment"
	//   },
	//   "response": {
	//     "$ref": "Comment"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.comments.list":

type CommentsListCall struct {
	s      *Service
	fileId string
	opt_   map[string]interface{}
}

// List: Lists a file's comments.
func (r *CommentsService) List(fileId string) *CommentsListCall {
	c := &CommentsListCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	return c
}

// IncludeDeleted sets the optional parameter "includeDeleted": If set,
// all comments and replies, including deleted comments and replies
// (with content stripped) will be returned.
func (c *CommentsListCall) IncludeDeleted(includeDeleted bool) *CommentsListCall {
	c.opt_["includeDeleted"] = includeDeleted
	return c
}

// MaxResults sets the optional parameter "maxResults": The maximum
// number of discussions to include in the response, used for paging.
func (c *CommentsListCall) MaxResults(maxResults int64) *CommentsListCall {
	c.opt_["maxResults"] = maxResults
	return c
}

// PageToken sets the optional parameter "pageToken": The continuation
// token, used to page through large result sets. To get the next page
// of results, set this parameter to the value of "nextPageToken" from
// the previous response.
func (c *CommentsListCall) PageToken(pageToken string) *CommentsListCall {
	c.opt_["pageToken"] = pageToken
	return c
}

// UpdatedMin sets the optional parameter "updatedMin": Only discussions
// that were updated after this timestamp will be returned. Formatted as
// an RFC 3339 timestamp.
func (c *CommentsListCall) UpdatedMin(updatedMin string) *CommentsListCall {
	c.opt_["updatedMin"] = updatedMin
	return c
}

func (c *CommentsListCall) Do() (*CommentList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["includeDeleted"]; ok {
		params.Set("includeDeleted", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["maxResults"]; ok {
		params.Set("maxResults", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["pageToken"]; ok {
		params.Set("pageToken", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updatedMin"]; ok {
		params.Set("updatedMin", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(CommentList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Lists a file's comments.",
	//   "httpMethod": "GET",
	//   "id": "drive.comments.list",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "includeDeleted": {
	//       "default": "false",
	//       "description": "If set, all comments and replies, including deleted comments and replies (with content stripped) will be returned.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "maxResults": {
	//       "default": "20",
	//       "description": "The maximum number of discussions to include in the response, used for paging.",
	//       "format": "int32",
	//       "location": "query",
	//       "maximum": "100",
	//       "minimum": "0",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "The continuation token, used to page through large result sets. To get the next page of results, set this parameter to the value of \"nextPageToken\" from the previous response.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "updatedMin": {
	//       "description": "Only discussions that were updated after this timestamp will be returned. Formatted as an RFC 3339 timestamp.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments",
	//   "response": {
	//     "$ref": "CommentList"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.comments.patch":

type CommentsPatchCall struct {
	s         *Service
	fileId    string
	commentId string
	comment   *Comment
	opt_      map[string]interface{}
}

// Patch: Updates an existing comment. This method supports patch
// semantics.
func (r *CommentsService) Patch(fileId string, commentId string, comment *Comment) *CommentsPatchCall {
	c := &CommentsPatchCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.commentId = commentId
	c.comment = comment
	return c
}

func (c *CommentsPatchCall) Do() (*Comment, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.comment)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments/{commentId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{commentId}", cleanPathString(c.commentId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Comment)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates an existing comment. This method supports patch semantics.",
	//   "httpMethod": "PATCH",
	//   "id": "drive.comments.patch",
	//   "parameterOrder": [
	//     "fileId",
	//     "commentId"
	//   ],
	//   "parameters": {
	//     "commentId": {
	//       "description": "The ID of the comment.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments/{commentId}",
	//   "request": {
	//     "$ref": "Comment"
	//   },
	//   "response": {
	//     "$ref": "Comment"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.comments.update":

type CommentsUpdateCall struct {
	s         *Service
	fileId    string
	commentId string
	comment   *Comment
	opt_      map[string]interface{}
}

// Update: Updates an existing comment.
func (r *CommentsService) Update(fileId string, commentId string, comment *Comment) *CommentsUpdateCall {
	c := &CommentsUpdateCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.commentId = commentId
	c.comment = comment
	return c
}

func (c *CommentsUpdateCall) Do() (*Comment, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.comment)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments/{commentId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{commentId}", cleanPathString(c.commentId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Comment)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates an existing comment.",
	//   "httpMethod": "PUT",
	//   "id": "drive.comments.update",
	//   "parameterOrder": [
	//     "fileId",
	//     "commentId"
	//   ],
	//   "parameters": {
	//     "commentId": {
	//       "description": "The ID of the comment.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments/{commentId}",
	//   "request": {
	//     "$ref": "Comment"
	//   },
	//   "response": {
	//     "$ref": "Comment"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.files.copy":

type FilesCopyCall struct {
	s      *Service
	fileId string
	file   *File
	opt_   map[string]interface{}
}

// Copy: Creates a copy of the specified file.
func (r *FilesService) Copy(fileId string, file *File) *FilesCopyCall {
	c := &FilesCopyCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.file = file
	return c
}

// Convert sets the optional parameter "convert": Whether to convert
// this file to the corresponding Google Docs format.
func (c *FilesCopyCall) Convert(convert bool) *FilesCopyCall {
	c.opt_["convert"] = convert
	return c
}

// Ocr sets the optional parameter "ocr": Whether to attempt OCR on
// .jpg, .png, .gif, or .pdf uploads.
func (c *FilesCopyCall) Ocr(ocr bool) *FilesCopyCall {
	c.opt_["ocr"] = ocr
	return c
}

// OcrLanguage sets the optional parameter "ocrLanguage": If ocr is
// true, hints at the language to use. Valid values are ISO 639-1 codes.
func (c *FilesCopyCall) OcrLanguage(ocrLanguage string) *FilesCopyCall {
	c.opt_["ocrLanguage"] = ocrLanguage
	return c
}

// Pinned sets the optional parameter "pinned": Whether to pin the head
// revision of the new copy.
func (c *FilesCopyCall) Pinned(pinned bool) *FilesCopyCall {
	c.opt_["pinned"] = pinned
	return c
}

// TimedTextLanguage sets the optional parameter "timedTextLanguage":
// The language of the timed text.
func (c *FilesCopyCall) TimedTextLanguage(timedTextLanguage string) *FilesCopyCall {
	c.opt_["timedTextLanguage"] = timedTextLanguage
	return c
}

// TimedTextTrackName sets the optional parameter "timedTextTrackName":
// The timed text track name.
func (c *FilesCopyCall) TimedTextTrackName(timedTextTrackName string) *FilesCopyCall {
	c.opt_["timedTextTrackName"] = timedTextTrackName
	return c
}

func (c *FilesCopyCall) Do() (*File, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.file)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["convert"]; ok {
		params.Set("convert", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["ocr"]; ok {
		params.Set("ocr", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["ocrLanguage"]; ok {
		params.Set("ocrLanguage", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["pinned"]; ok {
		params.Set("pinned", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["timedTextLanguage"]; ok {
		params.Set("timedTextLanguage", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["timedTextTrackName"]; ok {
		params.Set("timedTextTrackName", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/copy")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(File)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates a copy of the specified file.",
	//   "httpMethod": "POST",
	//   "id": "drive.files.copy",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "convert": {
	//       "default": "false",
	//       "description": "Whether to convert this file to the corresponding Google Docs format.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file to copy.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "ocr": {
	//       "default": "false",
	//       "description": "Whether to attempt OCR on .jpg, .png, .gif, or .pdf uploads.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "ocrLanguage": {
	//       "description": "If ocr is true, hints at the language to use. Valid values are ISO 639-1 codes.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pinned": {
	//       "default": "false",
	//       "description": "Whether to pin the head revision of the new copy.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "timedTextLanguage": {
	//       "description": "The language of the timed text.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "timedTextTrackName": {
	//       "description": "The timed text track name.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/copy",
	//   "request": {
	//     "$ref": "File"
	//   },
	//   "response": {
	//     "$ref": "File"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.files.delete":

type FilesDeleteCall struct {
	s      *Service
	fileId string
	opt_   map[string]interface{}
}

// Delete: Permanently deletes a file by ID. Skips the trash.
func (r *FilesService) Delete(fileId string) *FilesDeleteCall {
	c := &FilesDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	return c
}

func (c *FilesDeleteCall) Do() error {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
	// {
	//   "description": "Permanently deletes a file by ID. Skips the trash.",
	//   "httpMethod": "DELETE",
	//   "id": "drive.files.delete",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file to delete.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}",
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.files.get":

type FilesGetCall struct {
	s      *Service
	fileId string
	opt_   map[string]interface{}
}

// Get: Gets a file's metadata by ID.
func (r *FilesService) Get(fileId string) *FilesGetCall {
	c := &FilesGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	return c
}

// Projection sets the optional parameter "projection": This parameter
// is deprecated and has no function.
func (c *FilesGetCall) Projection(projection string) *FilesGetCall {
	c.opt_["projection"] = projection
	return c
}

// UpdateViewedDate sets the optional parameter "updateViewedDate":
// Whether to update the view date after successfully retrieving the
// file.
func (c *FilesGetCall) UpdateViewedDate(updateViewedDate bool) *FilesGetCall {
	c.opt_["updateViewedDate"] = updateViewedDate
	return c
}

func (c *FilesGetCall) Do() (*File, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["projection"]; ok {
		params.Set("projection", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updateViewedDate"]; ok {
		params.Set("updateViewedDate", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(File)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets a file's metadata by ID.",
	//   "httpMethod": "GET",
	//   "id": "drive.files.get",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID for the file in question.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "projection": {
	//       "description": "This parameter is deprecated and has no function.",
	//       "enum": [
	//         "BASIC",
	//         "FULL"
	//       ],
	//       "enumDescriptions": [
	//         "Deprecated",
	//         "Deprecated"
	//       ],
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "updateViewedDate": {
	//       "default": "false",
	//       "description": "Whether to update the view date after successfully retrieving the file.",
	//       "location": "query",
	//       "type": "boolean"
	//     }
	//   },
	//   "path": "files/{fileId}",
	//   "response": {
	//     "$ref": "File"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ],
	//   "supportsSubscription": true
	// }

}

// method id "drive.files.insert":

type FilesInsertCall struct {
	s      *Service
	file   *File
	opt_   map[string]interface{}
	media_ io.Reader
}

// Insert: Insert a new file.
func (r *FilesService) Insert(file *File) *FilesInsertCall {
	c := &FilesInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.file = file
	return c
}

// Convert sets the optional parameter "convert": Whether to convert
// this file to the corresponding Google Docs format.
func (c *FilesInsertCall) Convert(convert bool) *FilesInsertCall {
	c.opt_["convert"] = convert
	return c
}

// Ocr sets the optional parameter "ocr": Whether to attempt OCR on
// .jpg, .png, .gif, or .pdf uploads.
func (c *FilesInsertCall) Ocr(ocr bool) *FilesInsertCall {
	c.opt_["ocr"] = ocr
	return c
}

// OcrLanguage sets the optional parameter "ocrLanguage": If ocr is
// true, hints at the language to use. Valid values are ISO 639-1 codes.
func (c *FilesInsertCall) OcrLanguage(ocrLanguage string) *FilesInsertCall {
	c.opt_["ocrLanguage"] = ocrLanguage
	return c
}

// Pinned sets the optional parameter "pinned": Whether to pin the head
// revision of the uploaded file.
func (c *FilesInsertCall) Pinned(pinned bool) *FilesInsertCall {
	c.opt_["pinned"] = pinned
	return c
}

// TimedTextLanguage sets the optional parameter "timedTextLanguage":
// The language of the timed text.
func (c *FilesInsertCall) TimedTextLanguage(timedTextLanguage string) *FilesInsertCall {
	c.opt_["timedTextLanguage"] = timedTextLanguage
	return c
}

// TimedTextTrackName sets the optional parameter "timedTextTrackName":
// The timed text track name.
func (c *FilesInsertCall) TimedTextTrackName(timedTextTrackName string) *FilesInsertCall {
	c.opt_["timedTextTrackName"] = timedTextTrackName
	return c
}

// UseContentAsIndexableText sets the optional parameter
// "useContentAsIndexableText": Whether to use the content as indexable
// text.
func (c *FilesInsertCall) UseContentAsIndexableText(useContentAsIndexableText bool) *FilesInsertCall {
	c.opt_["useContentAsIndexableText"] = useContentAsIndexableText
	return c
}
func (c *FilesInsertCall) Media(r io.Reader) *FilesInsertCall {
	c.media_ = r
	return c
}

func (c *FilesInsertCall) Do() (*File, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.file)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["convert"]; ok {
		params.Set("convert", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["ocr"]; ok {
		params.Set("ocr", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["ocrLanguage"]; ok {
		params.Set("ocrLanguage", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["pinned"]; ok {
		params.Set("pinned", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["timedTextLanguage"]; ok {
		params.Set("timedTextLanguage", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["timedTextTrackName"]; ok {
		params.Set("timedTextTrackName", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["useContentAsIndexableText"]; ok {
		params.Set("useContentAsIndexableText", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files")
	if c.media_ != nil {
		urls = strings.Replace(urls, "https://www.googleapis.com/", "https://www.googleapis.com/upload/", 1)
		params.Set("uploadType", "multipart")
	}
	urls += "?" + params.Encode()
	contentLength_, hasMedia_ := googleapi.ConditionallyIncludeMedia(c.media_, &body, &ctype)
	req, _ := http.NewRequest("POST", urls, body)
	if hasMedia_ {
		req.ContentLength = contentLength_
	}
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(File)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Insert a new file.",
	//   "httpMethod": "POST",
	//   "id": "drive.files.insert",
	//   "mediaUpload": {
	//     "accept": [
	//       "*/*"
	//     ],
	//     "maxSize": "10GB",
	//     "protocols": {
	//       "resumable": {
	//         "multipart": true,
	//         "path": "/resumable/upload/drive/v2/files"
	//       },
	//       "simple": {
	//         "multipart": true,
	//         "path": "/upload/drive/v2/files"
	//       }
	//     }
	//   },
	//   "parameters": {
	//     "convert": {
	//       "default": "false",
	//       "description": "Whether to convert this file to the corresponding Google Docs format.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "ocr": {
	//       "default": "false",
	//       "description": "Whether to attempt OCR on .jpg, .png, .gif, or .pdf uploads.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "ocrLanguage": {
	//       "description": "If ocr is true, hints at the language to use. Valid values are ISO 639-1 codes.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pinned": {
	//       "default": "false",
	//       "description": "Whether to pin the head revision of the uploaded file.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "timedTextLanguage": {
	//       "description": "The language of the timed text.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "timedTextTrackName": {
	//       "description": "The timed text track name.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "useContentAsIndexableText": {
	//       "default": "false",
	//       "description": "Whether to use the content as indexable text.",
	//       "location": "query",
	//       "type": "boolean"
	//     }
	//   },
	//   "path": "files",
	//   "request": {
	//     "$ref": "File"
	//   },
	//   "response": {
	//     "$ref": "File"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ],
	//   "supportsMediaUpload": true,
	//   "supportsSubscription": true
	// }

}

// method id "drive.files.list":

type FilesListCall struct {
	s    *Service
	opt_ map[string]interface{}
}

// List: Lists the user's files.
func (r *FilesService) List() *FilesListCall {
	c := &FilesListCall{s: r.s, opt_: make(map[string]interface{})}
	return c
}

// MaxResults sets the optional parameter "maxResults": Maximum number
// of files to return.
func (c *FilesListCall) MaxResults(maxResults int64) *FilesListCall {
	c.opt_["maxResults"] = maxResults
	return c
}

// PageToken sets the optional parameter "pageToken": Page token for
// files.
func (c *FilesListCall) PageToken(pageToken string) *FilesListCall {
	c.opt_["pageToken"] = pageToken
	return c
}

// Projection sets the optional parameter "projection": This parameter
// is deprecated and has no function.
func (c *FilesListCall) Projection(projection string) *FilesListCall {
	c.opt_["projection"] = projection
	return c
}

// Q sets the optional parameter "q": Query string for searching files.
func (c *FilesListCall) Q(q string) *FilesListCall {
	c.opt_["q"] = q
	return c
}

func (c *FilesListCall) Do() (*FileList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["maxResults"]; ok {
		params.Set("maxResults", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["pageToken"]; ok {
		params.Set("pageToken", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["projection"]; ok {
		params.Set("projection", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["q"]; ok {
		params.Set("q", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(FileList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Lists the user's files.",
	//   "httpMethod": "GET",
	//   "id": "drive.files.list",
	//   "parameters": {
	//     "maxResults": {
	//       "default": "100",
	//       "description": "Maximum number of files to return.",
	//       "format": "int32",
	//       "location": "query",
	//       "minimum": "0",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "Page token for files.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "projection": {
	//       "description": "This parameter is deprecated and has no function.",
	//       "enum": [
	//         "BASIC",
	//         "FULL"
	//       ],
	//       "enumDescriptions": [
	//         "Deprecated",
	//         "Deprecated"
	//       ],
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "q": {
	//       "description": "Query string for searching files.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "files",
	//   "response": {
	//     "$ref": "FileList"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.files.patch":

type FilesPatchCall struct {
	s      *Service
	fileId string
	file   *File
	opt_   map[string]interface{}
}

// Patch: Updates file metadata and/or content. This method supports
// patch semantics.
func (r *FilesService) Patch(fileId string, file *File) *FilesPatchCall {
	c := &FilesPatchCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.file = file
	return c
}

// Convert sets the optional parameter "convert": Whether to convert
// this file to the corresponding Google Docs format.
func (c *FilesPatchCall) Convert(convert bool) *FilesPatchCall {
	c.opt_["convert"] = convert
	return c
}

// NewRevision sets the optional parameter "newRevision": Whether a blob
// upload should create a new revision. If not set or false, the blob
// data in the current head revision is replaced. If true, a new blob is
// created as head revision, and previous revisions are preserved
// (causing increased use of the user's data storage quota).
func (c *FilesPatchCall) NewRevision(newRevision bool) *FilesPatchCall {
	c.opt_["newRevision"] = newRevision
	return c
}

// Ocr sets the optional parameter "ocr": Whether to attempt OCR on
// .jpg, .png, .gif, or .pdf uploads.
func (c *FilesPatchCall) Ocr(ocr bool) *FilesPatchCall {
	c.opt_["ocr"] = ocr
	return c
}

// OcrLanguage sets the optional parameter "ocrLanguage": If ocr is
// true, hints at the language to use. Valid values are ISO 639-1 codes.
func (c *FilesPatchCall) OcrLanguage(ocrLanguage string) *FilesPatchCall {
	c.opt_["ocrLanguage"] = ocrLanguage
	return c
}

// Pinned sets the optional parameter "pinned": Whether to pin the new
// revision.
func (c *FilesPatchCall) Pinned(pinned bool) *FilesPatchCall {
	c.opt_["pinned"] = pinned
	return c
}

// SetModifiedDate sets the optional parameter "setModifiedDate":
// Whether to set the modified date with the supplied modified date.
func (c *FilesPatchCall) SetModifiedDate(setModifiedDate bool) *FilesPatchCall {
	c.opt_["setModifiedDate"] = setModifiedDate
	return c
}

// TimedTextLanguage sets the optional parameter "timedTextLanguage":
// The language of the timed text.
func (c *FilesPatchCall) TimedTextLanguage(timedTextLanguage string) *FilesPatchCall {
	c.opt_["timedTextLanguage"] = timedTextLanguage
	return c
}

// TimedTextTrackName sets the optional parameter "timedTextTrackName":
// The timed text track name.
func (c *FilesPatchCall) TimedTextTrackName(timedTextTrackName string) *FilesPatchCall {
	c.opt_["timedTextTrackName"] = timedTextTrackName
	return c
}

// UpdateViewedDate sets the optional parameter "updateViewedDate":
// Whether to update the view date after successfully updating the file.
func (c *FilesPatchCall) UpdateViewedDate(updateViewedDate bool) *FilesPatchCall {
	c.opt_["updateViewedDate"] = updateViewedDate
	return c
}

// UseContentAsIndexableText sets the optional parameter
// "useContentAsIndexableText": Whether to use the content as indexable
// text.
func (c *FilesPatchCall) UseContentAsIndexableText(useContentAsIndexableText bool) *FilesPatchCall {
	c.opt_["useContentAsIndexableText"] = useContentAsIndexableText
	return c
}

func (c *FilesPatchCall) Do() (*File, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.file)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["convert"]; ok {
		params.Set("convert", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["newRevision"]; ok {
		params.Set("newRevision", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["ocr"]; ok {
		params.Set("ocr", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["ocrLanguage"]; ok {
		params.Set("ocrLanguage", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["pinned"]; ok {
		params.Set("pinned", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["setModifiedDate"]; ok {
		params.Set("setModifiedDate", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["timedTextLanguage"]; ok {
		params.Set("timedTextLanguage", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["timedTextTrackName"]; ok {
		params.Set("timedTextTrackName", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updateViewedDate"]; ok {
		params.Set("updateViewedDate", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["useContentAsIndexableText"]; ok {
		params.Set("useContentAsIndexableText", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(File)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates file metadata and/or content. This method supports patch semantics.",
	//   "httpMethod": "PATCH",
	//   "id": "drive.files.patch",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "convert": {
	//       "default": "false",
	//       "description": "Whether to convert this file to the corresponding Google Docs format.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file to update.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "newRevision": {
	//       "default": "true",
	//       "description": "Whether a blob upload should create a new revision. If not set or false, the blob data in the current head revision is replaced. If true, a new blob is created as head revision, and previous revisions are preserved (causing increased use of the user's data storage quota).",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "ocr": {
	//       "default": "false",
	//       "description": "Whether to attempt OCR on .jpg, .png, .gif, or .pdf uploads.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "ocrLanguage": {
	//       "description": "If ocr is true, hints at the language to use. Valid values are ISO 639-1 codes.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pinned": {
	//       "default": "false",
	//       "description": "Whether to pin the new revision.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "setModifiedDate": {
	//       "default": "false",
	//       "description": "Whether to set the modified date with the supplied modified date.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "timedTextLanguage": {
	//       "description": "The language of the timed text.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "timedTextTrackName": {
	//       "description": "The timed text track name.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "updateViewedDate": {
	//       "default": "true",
	//       "description": "Whether to update the view date after successfully updating the file.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "useContentAsIndexableText": {
	//       "default": "false",
	//       "description": "Whether to use the content as indexable text.",
	//       "location": "query",
	//       "type": "boolean"
	//     }
	//   },
	//   "path": "files/{fileId}",
	//   "request": {
	//     "$ref": "File"
	//   },
	//   "response": {
	//     "$ref": "File"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.files.touch":

type FilesTouchCall struct {
	s      *Service
	fileId string
	opt_   map[string]interface{}
}

// Touch: Set the file's updated time to the current server time.
func (r *FilesService) Touch(fileId string) *FilesTouchCall {
	c := &FilesTouchCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	return c
}

func (c *FilesTouchCall) Do() (*File, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/touch")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(File)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Set the file's updated time to the current server time.",
	//   "httpMethod": "POST",
	//   "id": "drive.files.touch",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file to update.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/touch",
	//   "response": {
	//     "$ref": "File"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.files.trash":

type FilesTrashCall struct {
	s      *Service
	fileId string
	opt_   map[string]interface{}
}

// Trash: Moves a file to the trash.
func (r *FilesService) Trash(fileId string) *FilesTrashCall {
	c := &FilesTrashCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	return c
}

func (c *FilesTrashCall) Do() (*File, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/trash")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(File)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Moves a file to the trash.",
	//   "httpMethod": "POST",
	//   "id": "drive.files.trash",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file to trash.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/trash",
	//   "response": {
	//     "$ref": "File"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.files.untrash":

type FilesUntrashCall struct {
	s      *Service
	fileId string
	opt_   map[string]interface{}
}

// Untrash: Restores a file from the trash.
func (r *FilesService) Untrash(fileId string) *FilesUntrashCall {
	c := &FilesUntrashCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	return c
}

func (c *FilesUntrashCall) Do() (*File, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/untrash")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(File)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Restores a file from the trash.",
	//   "httpMethod": "POST",
	//   "id": "drive.files.untrash",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file to untrash.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/untrash",
	//   "response": {
	//     "$ref": "File"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.files.update":

type FilesUpdateCall struct {
	s      *Service
	fileId string
	file   *File
	opt_   map[string]interface{}
	media_ io.Reader
}

// Update: Updates file metadata and/or content.
func (r *FilesService) Update(fileId string, file *File) *FilesUpdateCall {
	c := &FilesUpdateCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.file = file
	return c
}

// Convert sets the optional parameter "convert": Whether to convert
// this file to the corresponding Google Docs format.
func (c *FilesUpdateCall) Convert(convert bool) *FilesUpdateCall {
	c.opt_["convert"] = convert
	return c
}

// NewRevision sets the optional parameter "newRevision": Whether a blob
// upload should create a new revision. If not set or false, the blob
// data in the current head revision is replaced. If true, a new blob is
// created as head revision, and previous revisions are preserved
// (causing increased use of the user's data storage quota).
func (c *FilesUpdateCall) NewRevision(newRevision bool) *FilesUpdateCall {
	c.opt_["newRevision"] = newRevision
	return c
}

// Ocr sets the optional parameter "ocr": Whether to attempt OCR on
// .jpg, .png, .gif, or .pdf uploads.
func (c *FilesUpdateCall) Ocr(ocr bool) *FilesUpdateCall {
	c.opt_["ocr"] = ocr
	return c
}

// OcrLanguage sets the optional parameter "ocrLanguage": If ocr is
// true, hints at the language to use. Valid values are ISO 639-1 codes.
func (c *FilesUpdateCall) OcrLanguage(ocrLanguage string) *FilesUpdateCall {
	c.opt_["ocrLanguage"] = ocrLanguage
	return c
}

// Pinned sets the optional parameter "pinned": Whether to pin the new
// revision.
func (c *FilesUpdateCall) Pinned(pinned bool) *FilesUpdateCall {
	c.opt_["pinned"] = pinned
	return c
}

// SetModifiedDate sets the optional parameter "setModifiedDate":
// Whether to set the modified date with the supplied modified date.
func (c *FilesUpdateCall) SetModifiedDate(setModifiedDate bool) *FilesUpdateCall {
	c.opt_["setModifiedDate"] = setModifiedDate
	return c
}

// TimedTextLanguage sets the optional parameter "timedTextLanguage":
// The language of the timed text.
func (c *FilesUpdateCall) TimedTextLanguage(timedTextLanguage string) *FilesUpdateCall {
	c.opt_["timedTextLanguage"] = timedTextLanguage
	return c
}

// TimedTextTrackName sets the optional parameter "timedTextTrackName":
// The timed text track name.
func (c *FilesUpdateCall) TimedTextTrackName(timedTextTrackName string) *FilesUpdateCall {
	c.opt_["timedTextTrackName"] = timedTextTrackName
	return c
}

// UpdateViewedDate sets the optional parameter "updateViewedDate":
// Whether to update the view date after successfully updating the file.
func (c *FilesUpdateCall) UpdateViewedDate(updateViewedDate bool) *FilesUpdateCall {
	c.opt_["updateViewedDate"] = updateViewedDate
	return c
}

// UseContentAsIndexableText sets the optional parameter
// "useContentAsIndexableText": Whether to use the content as indexable
// text.
func (c *FilesUpdateCall) UseContentAsIndexableText(useContentAsIndexableText bool) *FilesUpdateCall {
	c.opt_["useContentAsIndexableText"] = useContentAsIndexableText
	return c
}
func (c *FilesUpdateCall) Media(r io.Reader) *FilesUpdateCall {
	c.media_ = r
	return c
}

func (c *FilesUpdateCall) Do() (*File, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.file)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["convert"]; ok {
		params.Set("convert", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["newRevision"]; ok {
		params.Set("newRevision", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["ocr"]; ok {
		params.Set("ocr", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["ocrLanguage"]; ok {
		params.Set("ocrLanguage", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["pinned"]; ok {
		params.Set("pinned", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["setModifiedDate"]; ok {
		params.Set("setModifiedDate", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["timedTextLanguage"]; ok {
		params.Set("timedTextLanguage", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["timedTextTrackName"]; ok {
		params.Set("timedTextTrackName", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["updateViewedDate"]; ok {
		params.Set("updateViewedDate", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["useContentAsIndexableText"]; ok {
		params.Set("useContentAsIndexableText", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}")
	if c.media_ != nil {
		urls = strings.Replace(urls, "https://www.googleapis.com/", "https://www.googleapis.com/upload/", 1)
		params.Set("uploadType", "multipart")
	}
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	contentLength_, hasMedia_ := googleapi.ConditionallyIncludeMedia(c.media_, &body, &ctype)
	req, _ := http.NewRequest("PUT", urls, body)
	if hasMedia_ {
		req.ContentLength = contentLength_
	}
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(File)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates file metadata and/or content.",
	//   "httpMethod": "PUT",
	//   "id": "drive.files.update",
	//   "mediaUpload": {
	//     "accept": [
	//       "*/*"
	//     ],
	//     "maxSize": "10GB",
	//     "protocols": {
	//       "resumable": {
	//         "multipart": true,
	//         "path": "/resumable/upload/drive/v2/files/{fileId}"
	//       },
	//       "simple": {
	//         "multipart": true,
	//         "path": "/upload/drive/v2/files/{fileId}"
	//       }
	//     }
	//   },
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "convert": {
	//       "default": "false",
	//       "description": "Whether to convert this file to the corresponding Google Docs format.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file to update.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "newRevision": {
	//       "default": "true",
	//       "description": "Whether a blob upload should create a new revision. If not set or false, the blob data in the current head revision is replaced. If true, a new blob is created as head revision, and previous revisions are preserved (causing increased use of the user's data storage quota).",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "ocr": {
	//       "default": "false",
	//       "description": "Whether to attempt OCR on .jpg, .png, .gif, or .pdf uploads.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "ocrLanguage": {
	//       "description": "If ocr is true, hints at the language to use. Valid values are ISO 639-1 codes.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "pinned": {
	//       "default": "false",
	//       "description": "Whether to pin the new revision.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "setModifiedDate": {
	//       "default": "false",
	//       "description": "Whether to set the modified date with the supplied modified date.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "timedTextLanguage": {
	//       "description": "The language of the timed text.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "timedTextTrackName": {
	//       "description": "The timed text track name.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "updateViewedDate": {
	//       "default": "true",
	//       "description": "Whether to update the view date after successfully updating the file.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "useContentAsIndexableText": {
	//       "default": "false",
	//       "description": "Whether to use the content as indexable text.",
	//       "location": "query",
	//       "type": "boolean"
	//     }
	//   },
	//   "path": "files/{fileId}",
	//   "request": {
	//     "$ref": "File"
	//   },
	//   "response": {
	//     "$ref": "File"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ],
	//   "supportsMediaUpload": true
	// }

}

// method id "drive.parents.delete":

type ParentsDeleteCall struct {
	s        *Service
	fileId   string
	parentId string
	opt_     map[string]interface{}
}

// Delete: Removes a parent from a file.
func (r *ParentsService) Delete(fileId string, parentId string) *ParentsDeleteCall {
	c := &ParentsDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.parentId = parentId
	return c
}

func (c *ParentsDeleteCall) Do() error {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/parents/{parentId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{parentId}", cleanPathString(c.parentId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
	// {
	//   "description": "Removes a parent from a file.",
	//   "httpMethod": "DELETE",
	//   "id": "drive.parents.delete",
	//   "parameterOrder": [
	//     "fileId",
	//     "parentId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "parentId": {
	//       "description": "The ID of the parent.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/parents/{parentId}",
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.parents.get":

type ParentsGetCall struct {
	s        *Service
	fileId   string
	parentId string
	opt_     map[string]interface{}
}

// Get: Gets a specific parent reference.
func (r *ParentsService) Get(fileId string, parentId string) *ParentsGetCall {
	c := &ParentsGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.parentId = parentId
	return c
}

func (c *ParentsGetCall) Do() (*ParentReference, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/parents/{parentId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{parentId}", cleanPathString(c.parentId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(ParentReference)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets a specific parent reference.",
	//   "httpMethod": "GET",
	//   "id": "drive.parents.get",
	//   "parameterOrder": [
	//     "fileId",
	//     "parentId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "parentId": {
	//       "description": "The ID of the parent.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/parents/{parentId}",
	//   "response": {
	//     "$ref": "ParentReference"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.parents.insert":

type ParentsInsertCall struct {
	s               *Service
	fileId          string
	parentreference *ParentReference
	opt_            map[string]interface{}
}

// Insert: Adds a parent folder for a file.
func (r *ParentsService) Insert(fileId string, parentreference *ParentReference) *ParentsInsertCall {
	c := &ParentsInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.parentreference = parentreference
	return c
}

func (c *ParentsInsertCall) Do() (*ParentReference, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.parentreference)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/parents")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(ParentReference)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Adds a parent folder for a file.",
	//   "httpMethod": "POST",
	//   "id": "drive.parents.insert",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/parents",
	//   "request": {
	//     "$ref": "ParentReference"
	//   },
	//   "response": {
	//     "$ref": "ParentReference"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.parents.list":

type ParentsListCall struct {
	s      *Service
	fileId string
	opt_   map[string]interface{}
}

// List: Lists a file's parents.
func (r *ParentsService) List(fileId string) *ParentsListCall {
	c := &ParentsListCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	return c
}

func (c *ParentsListCall) Do() (*ParentList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/parents")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(ParentList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Lists a file's parents.",
	//   "httpMethod": "GET",
	//   "id": "drive.parents.list",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/parents",
	//   "response": {
	//     "$ref": "ParentList"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.permissions.delete":

type PermissionsDeleteCall struct {
	s            *Service
	fileId       string
	permissionId string
	opt_         map[string]interface{}
}

// Delete: Deletes a permission from a file.
func (r *PermissionsService) Delete(fileId string, permissionId string) *PermissionsDeleteCall {
	c := &PermissionsDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.permissionId = permissionId
	return c
}

func (c *PermissionsDeleteCall) Do() error {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/permissions/{permissionId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{permissionId}", cleanPathString(c.permissionId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
	// {
	//   "description": "Deletes a permission from a file.",
	//   "httpMethod": "DELETE",
	//   "id": "drive.permissions.delete",
	//   "parameterOrder": [
	//     "fileId",
	//     "permissionId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID for the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "permissionId": {
	//       "description": "The ID for the permission.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/permissions/{permissionId}",
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.permissions.get":

type PermissionsGetCall struct {
	s            *Service
	fileId       string
	permissionId string
	opt_         map[string]interface{}
}

// Get: Gets a permission by ID.
func (r *PermissionsService) Get(fileId string, permissionId string) *PermissionsGetCall {
	c := &PermissionsGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.permissionId = permissionId
	return c
}

func (c *PermissionsGetCall) Do() (*Permission, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/permissions/{permissionId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{permissionId}", cleanPathString(c.permissionId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Permission)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets a permission by ID.",
	//   "httpMethod": "GET",
	//   "id": "drive.permissions.get",
	//   "parameterOrder": [
	//     "fileId",
	//     "permissionId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID for the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "permissionId": {
	//       "description": "The ID for the permission.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/permissions/{permissionId}",
	//   "response": {
	//     "$ref": "Permission"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.permissions.insert":

type PermissionsInsertCall struct {
	s          *Service
	fileId     string
	permission *Permission
	opt_       map[string]interface{}
}

// Insert: Inserts a permission for a file.
func (r *PermissionsService) Insert(fileId string, permission *Permission) *PermissionsInsertCall {
	c := &PermissionsInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.permission = permission
	return c
}

// EmailMessage sets the optional parameter "emailMessage": A custom
// message to include in notification emails.
func (c *PermissionsInsertCall) EmailMessage(emailMessage string) *PermissionsInsertCall {
	c.opt_["emailMessage"] = emailMessage
	return c
}

// SendNotificationEmails sets the optional parameter
// "sendNotificationEmails": Whether to send notification emails when
// sharing to users or groups.
func (c *PermissionsInsertCall) SendNotificationEmails(sendNotificationEmails bool) *PermissionsInsertCall {
	c.opt_["sendNotificationEmails"] = sendNotificationEmails
	return c
}

func (c *PermissionsInsertCall) Do() (*Permission, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.permission)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["emailMessage"]; ok {
		params.Set("emailMessage", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["sendNotificationEmails"]; ok {
		params.Set("sendNotificationEmails", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/permissions")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Permission)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Inserts a permission for a file.",
	//   "httpMethod": "POST",
	//   "id": "drive.permissions.insert",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "emailMessage": {
	//       "description": "A custom message to include in notification emails.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID for the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "sendNotificationEmails": {
	//       "default": "true",
	//       "description": "Whether to send notification emails when sharing to users or groups.",
	//       "location": "query",
	//       "type": "boolean"
	//     }
	//   },
	//   "path": "files/{fileId}/permissions",
	//   "request": {
	//     "$ref": "Permission"
	//   },
	//   "response": {
	//     "$ref": "Permission"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.permissions.list":

type PermissionsListCall struct {
	s      *Service
	fileId string
	opt_   map[string]interface{}
}

// List: Lists a file's permissions.
func (r *PermissionsService) List(fileId string) *PermissionsListCall {
	c := &PermissionsListCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	return c
}

func (c *PermissionsListCall) Do() (*PermissionList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/permissions")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(PermissionList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Lists a file's permissions.",
	//   "httpMethod": "GET",
	//   "id": "drive.permissions.list",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID for the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/permissions",
	//   "response": {
	//     "$ref": "PermissionList"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.permissions.patch":

type PermissionsPatchCall struct {
	s            *Service
	fileId       string
	permissionId string
	permission   *Permission
	opt_         map[string]interface{}
}

// Patch: Updates a permission. This method supports patch semantics.
func (r *PermissionsService) Patch(fileId string, permissionId string, permission *Permission) *PermissionsPatchCall {
	c := &PermissionsPatchCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.permissionId = permissionId
	c.permission = permission
	return c
}

func (c *PermissionsPatchCall) Do() (*Permission, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.permission)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/permissions/{permissionId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{permissionId}", cleanPathString(c.permissionId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Permission)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates a permission. This method supports patch semantics.",
	//   "httpMethod": "PATCH",
	//   "id": "drive.permissions.patch",
	//   "parameterOrder": [
	//     "fileId",
	//     "permissionId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID for the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "permissionId": {
	//       "description": "The ID for the permission.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/permissions/{permissionId}",
	//   "request": {
	//     "$ref": "Permission"
	//   },
	//   "response": {
	//     "$ref": "Permission"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.permissions.update":

type PermissionsUpdateCall struct {
	s            *Service
	fileId       string
	permissionId string
	permission   *Permission
	opt_         map[string]interface{}
}

// Update: Updates a permission.
func (r *PermissionsService) Update(fileId string, permissionId string, permission *Permission) *PermissionsUpdateCall {
	c := &PermissionsUpdateCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.permissionId = permissionId
	c.permission = permission
	return c
}

func (c *PermissionsUpdateCall) Do() (*Permission, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.permission)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/permissions/{permissionId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{permissionId}", cleanPathString(c.permissionId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Permission)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates a permission.",
	//   "httpMethod": "PUT",
	//   "id": "drive.permissions.update",
	//   "parameterOrder": [
	//     "fileId",
	//     "permissionId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID for the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "permissionId": {
	//       "description": "The ID for the permission.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/permissions/{permissionId}",
	//   "request": {
	//     "$ref": "Permission"
	//   },
	//   "response": {
	//     "$ref": "Permission"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.replies.delete":

type RepliesDeleteCall struct {
	s         *Service
	fileId    string
	commentId string
	replyId   string
	opt_      map[string]interface{}
}

// Delete: Deletes a reply.
func (r *RepliesService) Delete(fileId string, commentId string, replyId string) *RepliesDeleteCall {
	c := &RepliesDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.commentId = commentId
	c.replyId = replyId
	return c
}

func (c *RepliesDeleteCall) Do() error {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments/{commentId}/replies/{replyId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{commentId}", cleanPathString(c.commentId), 1)
	urls = strings.Replace(urls, "{replyId}", cleanPathString(c.replyId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
	// {
	//   "description": "Deletes a reply.",
	//   "httpMethod": "DELETE",
	//   "id": "drive.replies.delete",
	//   "parameterOrder": [
	//     "fileId",
	//     "commentId",
	//     "replyId"
	//   ],
	//   "parameters": {
	//     "commentId": {
	//       "description": "The ID of the comment.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "replyId": {
	//       "description": "The ID of the reply.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments/{commentId}/replies/{replyId}",
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.replies.get":

type RepliesGetCall struct {
	s         *Service
	fileId    string
	commentId string
	replyId   string
	opt_      map[string]interface{}
}

// Get: Gets a reply.
func (r *RepliesService) Get(fileId string, commentId string, replyId string) *RepliesGetCall {
	c := &RepliesGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.commentId = commentId
	c.replyId = replyId
	return c
}

// IncludeDeleted sets the optional parameter "includeDeleted": If set,
// this will succeed when retrieving a deleted reply.
func (c *RepliesGetCall) IncludeDeleted(includeDeleted bool) *RepliesGetCall {
	c.opt_["includeDeleted"] = includeDeleted
	return c
}

func (c *RepliesGetCall) Do() (*CommentReply, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["includeDeleted"]; ok {
		params.Set("includeDeleted", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments/{commentId}/replies/{replyId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{commentId}", cleanPathString(c.commentId), 1)
	urls = strings.Replace(urls, "{replyId}", cleanPathString(c.replyId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(CommentReply)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets a reply.",
	//   "httpMethod": "GET",
	//   "id": "drive.replies.get",
	//   "parameterOrder": [
	//     "fileId",
	//     "commentId",
	//     "replyId"
	//   ],
	//   "parameters": {
	//     "commentId": {
	//       "description": "The ID of the comment.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "includeDeleted": {
	//       "default": "false",
	//       "description": "If set, this will succeed when retrieving a deleted reply.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "replyId": {
	//       "description": "The ID of the reply.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments/{commentId}/replies/{replyId}",
	//   "response": {
	//     "$ref": "CommentReply"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.replies.insert":

type RepliesInsertCall struct {
	s            *Service
	fileId       string
	commentId    string
	commentreply *CommentReply
	opt_         map[string]interface{}
}

// Insert: Creates a new reply to the given comment.
func (r *RepliesService) Insert(fileId string, commentId string, commentreply *CommentReply) *RepliesInsertCall {
	c := &RepliesInsertCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.commentId = commentId
	c.commentreply = commentreply
	return c
}

func (c *RepliesInsertCall) Do() (*CommentReply, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.commentreply)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments/{commentId}/replies")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{commentId}", cleanPathString(c.commentId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(CommentReply)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Creates a new reply to the given comment.",
	//   "httpMethod": "POST",
	//   "id": "drive.replies.insert",
	//   "parameterOrder": [
	//     "fileId",
	//     "commentId"
	//   ],
	//   "parameters": {
	//     "commentId": {
	//       "description": "The ID of the comment.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments/{commentId}/replies",
	//   "request": {
	//     "$ref": "CommentReply"
	//   },
	//   "response": {
	//     "$ref": "CommentReply"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.replies.list":

type RepliesListCall struct {
	s         *Service
	fileId    string
	commentId string
	opt_      map[string]interface{}
}

// List: Lists all of the replies to a comment.
func (r *RepliesService) List(fileId string, commentId string) *RepliesListCall {
	c := &RepliesListCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.commentId = commentId
	return c
}

// IncludeDeleted sets the optional parameter "includeDeleted": If set,
// all replies, including deleted replies (with content stripped) will
// be returned.
func (c *RepliesListCall) IncludeDeleted(includeDeleted bool) *RepliesListCall {
	c.opt_["includeDeleted"] = includeDeleted
	return c
}

// MaxResults sets the optional parameter "maxResults": The maximum
// number of replies to include in the response, used for paging.
func (c *RepliesListCall) MaxResults(maxResults int64) *RepliesListCall {
	c.opt_["maxResults"] = maxResults
	return c
}

// PageToken sets the optional parameter "pageToken": The continuation
// token, used to page through large result sets. To get the next page
// of results, set this parameter to the value of "nextPageToken" from
// the previous response.
func (c *RepliesListCall) PageToken(pageToken string) *RepliesListCall {
	c.opt_["pageToken"] = pageToken
	return c
}

func (c *RepliesListCall) Do() (*CommentReplyList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	if v, ok := c.opt_["includeDeleted"]; ok {
		params.Set("includeDeleted", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["maxResults"]; ok {
		params.Set("maxResults", fmt.Sprintf("%v", v))
	}
	if v, ok := c.opt_["pageToken"]; ok {
		params.Set("pageToken", fmt.Sprintf("%v", v))
	}
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments/{commentId}/replies")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{commentId}", cleanPathString(c.commentId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(CommentReplyList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Lists all of the replies to a comment.",
	//   "httpMethod": "GET",
	//   "id": "drive.replies.list",
	//   "parameterOrder": [
	//     "fileId",
	//     "commentId"
	//   ],
	//   "parameters": {
	//     "commentId": {
	//       "description": "The ID of the comment.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "includeDeleted": {
	//       "default": "false",
	//       "description": "If set, all replies, including deleted replies (with content stripped) will be returned.",
	//       "location": "query",
	//       "type": "boolean"
	//     },
	//     "maxResults": {
	//       "default": "20",
	//       "description": "The maximum number of replies to include in the response, used for paging.",
	//       "format": "int32",
	//       "location": "query",
	//       "maximum": "100",
	//       "minimum": "0",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "The continuation token, used to page through large result sets. To get the next page of results, set this parameter to the value of \"nextPageToken\" from the previous response.",
	//       "location": "query",
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments/{commentId}/replies",
	//   "response": {
	//     "$ref": "CommentReplyList"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.replies.patch":

type RepliesPatchCall struct {
	s            *Service
	fileId       string
	commentId    string
	replyId      string
	commentreply *CommentReply
	opt_         map[string]interface{}
}

// Patch: Updates an existing reply. This method supports patch
// semantics.
func (r *RepliesService) Patch(fileId string, commentId string, replyId string, commentreply *CommentReply) *RepliesPatchCall {
	c := &RepliesPatchCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.commentId = commentId
	c.replyId = replyId
	c.commentreply = commentreply
	return c
}

func (c *RepliesPatchCall) Do() (*CommentReply, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.commentreply)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments/{commentId}/replies/{replyId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{commentId}", cleanPathString(c.commentId), 1)
	urls = strings.Replace(urls, "{replyId}", cleanPathString(c.replyId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(CommentReply)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates an existing reply. This method supports patch semantics.",
	//   "httpMethod": "PATCH",
	//   "id": "drive.replies.patch",
	//   "parameterOrder": [
	//     "fileId",
	//     "commentId",
	//     "replyId"
	//   ],
	//   "parameters": {
	//     "commentId": {
	//       "description": "The ID of the comment.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "replyId": {
	//       "description": "The ID of the reply.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments/{commentId}/replies/{replyId}",
	//   "request": {
	//     "$ref": "CommentReply"
	//   },
	//   "response": {
	//     "$ref": "CommentReply"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.replies.update":

type RepliesUpdateCall struct {
	s            *Service
	fileId       string
	commentId    string
	replyId      string
	commentreply *CommentReply
	opt_         map[string]interface{}
}

// Update: Updates an existing reply.
func (r *RepliesService) Update(fileId string, commentId string, replyId string, commentreply *CommentReply) *RepliesUpdateCall {
	c := &RepliesUpdateCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.commentId = commentId
	c.replyId = replyId
	c.commentreply = commentreply
	return c
}

func (c *RepliesUpdateCall) Do() (*CommentReply, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.commentreply)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/comments/{commentId}/replies/{replyId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{commentId}", cleanPathString(c.commentId), 1)
	urls = strings.Replace(urls, "{replyId}", cleanPathString(c.replyId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(CommentReply)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates an existing reply.",
	//   "httpMethod": "PUT",
	//   "id": "drive.replies.update",
	//   "parameterOrder": [
	//     "fileId",
	//     "commentId",
	//     "replyId"
	//   ],
	//   "parameters": {
	//     "commentId": {
	//       "description": "The ID of the comment.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "replyId": {
	//       "description": "The ID of the reply.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/comments/{commentId}/replies/{replyId}",
	//   "request": {
	//     "$ref": "CommentReply"
	//   },
	//   "response": {
	//     "$ref": "CommentReply"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.revisions.delete":

type RevisionsDeleteCall struct {
	s          *Service
	fileId     string
	revisionId string
	opt_       map[string]interface{}
}

// Delete: Removes a revision.
func (r *RevisionsService) Delete(fileId string, revisionId string) *RevisionsDeleteCall {
	c := &RevisionsDeleteCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.revisionId = revisionId
	return c
}

func (c *RevisionsDeleteCall) Do() error {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/revisions/{revisionId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{revisionId}", cleanPathString(c.revisionId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("DELETE", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
	// {
	//   "description": "Removes a revision.",
	//   "httpMethod": "DELETE",
	//   "id": "drive.revisions.delete",
	//   "parameterOrder": [
	//     "fileId",
	//     "revisionId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "revisionId": {
	//       "description": "The ID of the revision.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/revisions/{revisionId}",
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.revisions.get":

type RevisionsGetCall struct {
	s          *Service
	fileId     string
	revisionId string
	opt_       map[string]interface{}
}

// Get: Gets a specific revision.
func (r *RevisionsService) Get(fileId string, revisionId string) *RevisionsGetCall {
	c := &RevisionsGetCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.revisionId = revisionId
	return c
}

func (c *RevisionsGetCall) Do() (*Revision, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/revisions/{revisionId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{revisionId}", cleanPathString(c.revisionId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Revision)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets a specific revision.",
	//   "httpMethod": "GET",
	//   "id": "drive.revisions.get",
	//   "parameterOrder": [
	//     "fileId",
	//     "revisionId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "revisionId": {
	//       "description": "The ID of the revision.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/revisions/{revisionId}",
	//   "response": {
	//     "$ref": "Revision"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.revisions.list":

type RevisionsListCall struct {
	s      *Service
	fileId string
	opt_   map[string]interface{}
}

// List: Lists a file's revisions.
func (r *RevisionsService) List(fileId string) *RevisionsListCall {
	c := &RevisionsListCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	return c
}

func (c *RevisionsListCall) Do() (*RevisionList, error) {
	var body io.Reader = nil
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/revisions")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("GET", urls, body)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(RevisionList)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Lists a file's revisions.",
	//   "httpMethod": "GET",
	//   "id": "drive.revisions.list",
	//   "parameterOrder": [
	//     "fileId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID of the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/revisions",
	//   "response": {
	//     "$ref": "RevisionList"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file",
	//     "https://www.googleapis.com/auth/drive.metadata.readonly",
	//     "https://www.googleapis.com/auth/drive.readonly"
	//   ]
	// }

}

// method id "drive.revisions.patch":

type RevisionsPatchCall struct {
	s          *Service
	fileId     string
	revisionId string
	revision   *Revision
	opt_       map[string]interface{}
}

// Patch: Updates a revision. This method supports patch semantics.
func (r *RevisionsService) Patch(fileId string, revisionId string, revision *Revision) *RevisionsPatchCall {
	c := &RevisionsPatchCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.revisionId = revisionId
	c.revision = revision
	return c
}

func (c *RevisionsPatchCall) Do() (*Revision, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.revision)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/revisions/{revisionId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{revisionId}", cleanPathString(c.revisionId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PATCH", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Revision)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates a revision. This method supports patch semantics.",
	//   "httpMethod": "PATCH",
	//   "id": "drive.revisions.patch",
	//   "parameterOrder": [
	//     "fileId",
	//     "revisionId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID for the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "revisionId": {
	//       "description": "The ID for the revision.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/revisions/{revisionId}",
	//   "request": {
	//     "$ref": "Revision"
	//   },
	//   "response": {
	//     "$ref": "Revision"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

// method id "drive.revisions.update":

type RevisionsUpdateCall struct {
	s          *Service
	fileId     string
	revisionId string
	revision   *Revision
	opt_       map[string]interface{}
}

// Update: Updates a revision.
func (r *RevisionsService) Update(fileId string, revisionId string, revision *Revision) *RevisionsUpdateCall {
	c := &RevisionsUpdateCall{s: r.s, opt_: make(map[string]interface{})}
	c.fileId = fileId
	c.revisionId = revisionId
	c.revision = revision
	return c
}

func (c *RevisionsUpdateCall) Do() (*Revision, error) {
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.revision)
	if err != nil {
		return nil, err
	}
	ctype := "application/json"
	params := make(url.Values)
	params.Set("alt", "json")
	urls := googleapi.ResolveRelative("https://www.googleapis.com/drive/v2/", "files/{fileId}/revisions/{revisionId}")
	urls = strings.Replace(urls, "{fileId}", cleanPathString(c.fileId), 1)
	urls = strings.Replace(urls, "{revisionId}", cleanPathString(c.revisionId), 1)
	urls += "?" + params.Encode()
	req, _ := http.NewRequest("PUT", urls, body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("User-Agent", "google-api-go-client/0.5")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Revision)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Updates a revision.",
	//   "httpMethod": "PUT",
	//   "id": "drive.revisions.update",
	//   "parameterOrder": [
	//     "fileId",
	//     "revisionId"
	//   ],
	//   "parameters": {
	//     "fileId": {
	//       "description": "The ID for the file.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "revisionId": {
	//       "description": "The ID for the revision.",
	//       "location": "path",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "files/{fileId}/revisions/{revisionId}",
	//   "request": {
	//     "$ref": "Revision"
	//   },
	//   "response": {
	//     "$ref": "Revision"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/drive",
	//     "https://www.googleapis.com/auth/drive.file"
	//   ]
	// }

}

func cleanPathString(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 0x2d && r <= 0x7a || r == '~' {
			return r
		}
		return -1
	}, s)
}
