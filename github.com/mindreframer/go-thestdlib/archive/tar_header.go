type Header struct {
    Name       string    // name of header file entry
    Mode       int64     // permission and mode bits
    Uid        int       // user id of owner
    Gid        int       // group id of owner
    Size       int64     // length in bytes
    ModTime    time.Time // modified time
    Typeflag   byte      // type of header entry
    Linkname   string    // target name of link
    Uname      string    // user name of owner
    Gname      string    // group name of owner
    Devmajor   int64     // major number of character or block device
    Devminor   int64     // minor number of character or block device
    AccessTime time.Time // access time
    ChangeTime time.Time // status change time
}
