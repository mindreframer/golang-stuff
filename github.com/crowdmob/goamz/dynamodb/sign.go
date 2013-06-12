package dynamodb

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const iSO8601BasicFormat = "20060102T150405Z"
const iSO8601BasicFormatShort = "20060102"

var (
	ErrNoDate = errors.New("Date header not supplied")
)

var lf = []byte{'\n'}

// Service represents an AWS-compatible service.
type Service struct {
	// Name is the name of the service being used (i.e. iam, etc)
	Name string

	// Region is the region you want to communicate with the service through. (i.e. us-east-1)
	Region string
}

func (s *Service) sign(k *aws.Auth, t time.Time) []byte {
	h := ghmac([]byte("AWS4"+k.SecretKey), []byte(t.Format(iSO8601BasicFormatShort)))
	h = ghmac(h, []byte(s.Region))
	h = ghmac(h, []byte(s.Name))
	h = ghmac(h, []byte("aws4_request"))
	return h
}

// For Testing.
func (s *Service) DerivedKey(k *aws.Auth, t time.Time) []byte {
	return s.sign(k, t)
}

// Sign signs an HTTP request with the given AWS keys for use on service s.
func (s *Service) Sign(keys *aws.Auth, r *http.Request) error {
	var t time.Time

	date := r.Header.Get("Date")
	if date == "" {
		return ErrNoDate
	}

	t, err := time.Parse(http.TimeFormat, date)
	if err != nil {
		return err
	}

	r.Header.Set("Date", t.Format(iSO8601BasicFormat))

	k := s.sign(keys, t)
	h := hmac.New(sha256.New, k)
	s.writeStringToSign(h, t, r)

	auth := bytes.NewBufferString("AWS4-HMAC-SHA256 ")
	auth.Write([]byte("Credential=" + keys.AccessKey + "/" + s.creds(t)))
	auth.Write([]byte{',', ' '})
	auth.Write([]byte("SignedHeaders="))
	s.writeHeaderList(auth, r)
	auth.Write([]byte{',', ' '})
	auth.Write([]byte("Signature=" + fmt.Sprintf("%x", h.Sum(nil))))

	r.Header.Set("Authorization", auth.String())

	return nil
}

func (s *Service) writeQuery(w io.Writer, r *http.Request) {
	var a []string
	for k, vs := range r.URL.Query() {
		k = url.QueryEscape(k)
		for _, v := range vs {
			if v == "" {
				a = append(a, k)
			} else {
				v = url.QueryEscape(v)
				a = append(a, k+"="+v)
			}
		}
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			w.Write([]byte{'&'})
		}

		w.Write([]byte(s))
	}
}

func (s *Service) writeHeader(w io.Writer, r *http.Request) {
	i, a := 0, make([]string, len(r.Header))
	for k, v := range r.Header {
		sort.Strings(v)
		a[i] = strings.ToLower(k) + ":" + strings.Join(v, ",")
		i++
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			w.Write(lf)
		}

		io.WriteString(w, s)
	}
}

func (s *Service) writeHeaderList(w io.Writer, r *http.Request) {
	i, a := 0, make([]string, len(r.Header))
	for k, _ := range r.Header {
		a[i] = strings.ToLower(k)
		i++
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			w.Write([]byte{';'})
		}
		w.Write([]byte(s))
	}

}

func (s *Service) writeBody(w io.Writer, r *http.Request) {
	var b []byte
	if r.Body == nil {
		b = []byte("")
	} else {
		var err error
		b, err = ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	h := sha256.New()
	h.Write(b)

	sum := h.Sum(nil)

	fmt.Fprintf(w, "%x", sum)
}

func (s *Service) writeURI(w io.Writer, r *http.Request) {
	path := r.URL.RequestURI()
	if r.URL.RawQuery != "" {
		path = path[:len(path)-len(r.URL.RawQuery)-1]
	}
	slash := strings.HasSuffix(path, "/")
	path = filepath.Clean(path)
	if path != "/" && slash {
		path += "/"
	}

	w.Write([]byte(path))
}

func (s *Service) writeRequest(w io.Writer, r *http.Request) {
	r.Header.Set("host", r.Host)

	w.Write([]byte(r.Method))
	w.Write(lf)
	s.writeURI(w, r)
	w.Write(lf)
	s.writeQuery(w, r)
	w.Write(lf)
	s.writeHeader(w, r)
	w.Write(lf)
	w.Write(lf)
	s.writeHeaderList(w, r)
	w.Write(lf)
	s.writeBody(w, r)
}

func (s *Service) writeStringToSign(w io.Writer, t time.Time, r *http.Request) {
	w.Write([]byte("AWS4-HMAC-SHA256"))
	w.Write(lf)
	w.Write([]byte(t.Format(iSO8601BasicFormat)))
	w.Write(lf)

	w.Write([]byte(s.creds(t)))
	w.Write(lf)

	h := sha256.New()
	s.writeRequest(h, r)
	fmt.Fprintf(w, "%x", h.Sum(nil))
}

func (s *Service) creds(t time.Time) string {
	return t.Format(iSO8601BasicFormatShort) + "/" + s.Region + "/" + s.Name + "/aws4_request"
}

func ghmac(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
