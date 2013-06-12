type FileHeader struct {
    Name             string
    CreatorVersion   uint16
    ReaderVersion    uint16
    Flags            uint16
    Method           uint16
    ModifiedTime     uint16 // MS-DOS time
    ModifiedDate     uint16 // MS-DOS date
    CRC32            uint32
    CompressedSize   uint32
    UncompressedSize uint32
    Extra            []byte
    ExternalAttrs    uint32 // Meaning depends on CreatorVersion
    Comment          string
}
