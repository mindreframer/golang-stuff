// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package firehose

import (
	"errors"
	"net/url"
	"time"
)

// The firehose event format implemented here is documented in:
//	https://github.com/tumblr/parmesan/wiki/Parmesan-API

// Event is a parsed representation of a Firehose event
type Event struct {

	// Activity describes the type of event (post create, liked, etc.)
	Activity Activity

	// Private equals true of this is a private event
	Private bool

	// PrivateData contains the contents of private structures included with the event
	PrivateData map[string]interface{}

	// Timestamp records the time the event occurred
	Timestamp time.Time

	// If the event Activity is CratePost, UpdatePost or DeletePost, Post holds event-specific information
	Post *Post

	// If the event Activity is Like or Unlike, Like holds event-specific information
	Like *Like
}

// Post keeps event information pertaining to CreatePost, UpdatePost and DeletePost events
type Post struct {
	ID          int64    // ID equals the post ID
	BlogID      int64    // BlogID equals the tumblelog ID of the owner
	BlogName    string   // BlogName is a textual representation of the tumblelog's identity
	PostURL     string   // PostURL is the public URL where the post can be viewed
	BlogURL     string   // BlogURL is the public URL of the main page of the owning tumblelog
	Type        PostType // Type determines if the post is text, photo, video, audio, etc.
	Tags        []string // Tags lists any user-supplied tags that apply to this post
	Title       string   // Title is the title of the post
	Body        string   // Body contains the post's HTML body
	Caption     string   // Caption is the caption of the post, when it applies (audio, video, etc.)
	SourceURL   string
	SourceTitle string
	Quote       string // For quote posts, Quote holds the contents of the quote
	LinkURL     string // LinkURL equals the link of this post, if post is of type link
	Description string
	Photos      []Photo // For photoset posts, Photos contains further photoset-specific details
}

// Photo represents a single photo inside a photoset
type Photo struct {
	Caption string     // Caption is the caption of the this photo
	Alt     []AltPhoto // Alt is slice of variations on the format of this photo
}

// BigestAlt returns a pointer to the largest format altenrative for this photo
func (ph *Photo) BiggestAlt() *AltPhoto {
	var alt *AltPhoto
	var width int
	for i, _ := range ph.Alt {
		x := &ph.Alt[i]
		if alt == nil || x.Width > width {
			alt, width = x, x.Width
		}
	}
	return alt
}

// AltPhoto represents a specific rendition of a photo in way of scaling
type AltPhoto struct {
	Width  int    // Width is the photo width in pixels
	Height int    // Height is the photo height in pixels
	URL    string // URL is the URL where the photo can be accessed
}

// Like holds details specific to like and unlike events
type Like struct {
	DestPostID   int64 // DestPostID is the post ID of the post that is being liked
	DestBlogID   int64 // DestBlogID is the tumblelog ID of the tumblelog owning the post that is being liked
	SourceBlogID int64 // SourceBlogID is the primary tumblelog ID of the user liking the post
	RootPostID   int64 // RootPostID is the post ID of the original post, if the liked post is a reblog
	RootBlogID   int64 // RootBlogID is the tumblelog ID of the tumblelog owning the original post, if the liked post is a reblog
	ParentPostID int64
	ParentBlogID int64
}

// The activity constants list the currently supported event activity types.
const (
	CreatePost Activity = iota
	UpdatePost
	DeletePost
	Likes
	Unlikes
	FirehoseCheckpoint
)

// Post constants enumerate all supported tumblelog post types.
type PostType byte

const (
	PostUnknown PostType = iota
	PostText
	PostQuote
	PostLink
	PostAnswer
	PostVideo
	PostAudio
	PostPhoto
	PostChat
)

// String returns a textual representation of the post type
func (pt PostType) String() string {
	switch pt {
	case PostText:
		return "text"
	case PostQuote:
		return "quote"
	case PostLink:
		return "link"
	case PostAnswer:
		return "answer"
	case PostVideo:
		return "video"
	case PostAudio:
		return "audio"
	case PostPhoto:
		return "photo"
	case PostChat:
		return "chat"
	}
	return "unknown"
}

func parsePostType(t string) PostType {
	switch t {
	case "text":
		return PostText
	case "quote":
		return PostQuote
	case "link":
		return PostLink
	case "answer":
		return PostAnswer
	case "video":
		return PostVideo
	case "audio":
		return PostAudio
	case "photo":
		return PostPhoto
	case "chat":
		return PostChat
	}
	return PostUnknown
}

// Errors returned by the parser of incoming events from the Tumblr Firehose
var (
	ErrParse   = errors.New("unrecognized semantics")
	ErrMissing = errors.New("missing field")
	ErrType    = errors.New("wrong type")
)

// IsSyntaxError returns true if err represents an event parsing error
func IsSyntaxError(err error) bool {
	switch err {
	case ErrParse, ErrMissing, ErrType:
		return true
	}
	return false
}

// Activity represents the type of user activity that an event holds
type Activity byte

// String returns a textual representation of the activity type
func (a Activity) String() string {
	switch a {
	case CreatePost:
		return "CreatePost"
	case UpdatePost:
		return "UpdatePost"
	case DeletePost:
		return "DeletePost"
	case Likes:
		return "Likes"
	case Unlikes:
		return "Unlikes"
	case FirehoseCheckpoint:
		return "FirehoseCheckpoint"
	}
	return "Unknown"
}

// ParseActivity converts a textual representation of an activity type to a compact representation of type Activity
func ParseActivity(activity string) (Activity, error) {
	switch activity {
	case "CreatePost":
		return CreatePost, nil
	case "UpdatePost":
		return UpdatePost, nil
	case "DeletePost":
		return DeletePost, nil
	case "Likes":
		return Likes, nil
	case "Unlikes":
		return Unlikes, nil
	case "FirehoseCheckpoint":
		return FirehoseCheckpoint, nil
	}
	return 0, ErrParse
}

func parseEvent(m map[string]interface{}) (ev *Event, err error) {
	/*
		defer func() {
			if err != nil {
				fmt.Printf("RAW:%#v\n", m)
			}
		}()
	*/

	ev = &Event{}
	if ev.Activity, err = ParseActivity(getString(m, "activity")); err != nil {
		return nil, err
	}
	if ev.Activity == FirehoseCheckpoint {
		return ev, nil
	}
	ev.Private = getBool(m, "isPrivate")
	ev.PrivateData = getMap(m, "privateJson")
	if epochms, err := getInt64(m, "timestamp"); err != nil {
		return nil, err
	} else {
		ev.Timestamp = time.Unix(0, epochms*1e6)
	}
	switch ev.Activity {
	case CreatePost, UpdatePost, DeletePost:
		p := &Post{}
		if p.ID, err = getInt64(m, "id"); err != nil {
			return nil, err
		}
		if p.BlogID, err = getInt64(m, "blogId"); err != nil {
			return nil, err
		}

		data := getMap(m, "data")
		if data == nil {
			return nil, ErrMissing
		}

		p.BlogName = getString(data, "blog_name")
		p.Type = parsePostType(getString(data, "type"))
		p.Title = getString(data, "title")
		p.Body = getString(data, "body")
		p.Caption = getString(data, "caption")
		p.SourceURL = getString(data, "source_url")
		p.SourceTitle = getString(data, "source_title")
		p.Quote = getString(data, "text")
		p.PostURL = getString(data, "post_url")
		p.BlogURL = blogFromPostURL(p.PostURL)
		p.LinkURL = getString(data, "url")
		p.Description = getString(data, "description")

		switch p.Type {
		case PostPhoto:
			p.Photos = parsePhotos(data)
		}

		// Tags
		if tags := getSlice(data, "tags"); tags != nil {
			for _, q := range tags {
				s, ok := q.(string)
				if ok {
					p.Tags = append(p.Tags, s)
				}
			}
		}

		ev.Post = p
		return ev, nil
	case Likes, Unlikes:
		l := &Like{}
		if l.DestPostID, err = getInt64(m, "destPostId"); err != nil {
			return nil, err
		}
		if l.DestBlogID, err = getInt64(m, "destBlogId"); err != nil {
			return nil, err
		}
		if l.SourceBlogID, err = getInt64(m, "sourceBlogId"); err != nil {
			return nil, err
		}
		if l.RootPostID, err = getInt64(m, "rootPostId"); err != nil {
			return nil, err
		}
		if l.RootBlogID, err = getInt64(m, "rootBlogId"); err != nil {
			return nil, err
		}
		if l.ParentPostID, err = getInt64(m, "parentPostId"); err != nil {
			return nil, err
		}
		if l.ParentBlogID, err = getInt64(m, "parentBlogId"); err != nil {
			return nil, err
		}
		ev.Like = l
		return ev, nil
	}
	return nil, ErrParse
}

func parsePhotos(data map[string]interface{}) []Photo {
	rawPhotos := getSlice(data, "photos")
	if rawPhotos == nil {
		return nil
	}
	photos := make([]Photo, len(rawPhotos))
	for i, rawPhoto_ := range rawPhotos {
		rawPhoto, ok := rawPhoto_.(map[string]interface{})
		if !ok {
			continue
		}
		var photo *Photo = &photos[i]
		photo.Caption = getString(rawPhoto, "caption")
		rawAlts := getSlice(rawPhoto, "alt_sizes")
		photo.Alt = make([]AltPhoto, len(rawAlts))
		for j, rawAlt_ := range rawAlts {
			rawAlt, ok := rawAlt_.(map[string]interface{})
			if !ok {
				continue
			}
			photo.Alt[j] = AltPhoto{
				Width:  getInt(rawAlt, "width"),
				Height: getInt(rawAlt, "height"),
				URL:    getString(rawAlt, "url"),
			}
		}

	}
	return photos
}

func blogFromPostURL(postURL string) string {
	u, err := url.ParseRequestURI(postURL)
	if err != nil {
		return ""
	}
	return u.Scheme + "://" + u.Host
}
