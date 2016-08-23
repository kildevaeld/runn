package filecrypt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/nacl/secretbox"
)

func read(r io.Reader, out []byte, nr int) error {
	re := 0
	try := 1
	for {
		log.Printf("reading try: %d", try)
		n, err := r.Read(out[re:])
		if err != nil {
			if n != nr {
				return errors.New("short read")
			}
			return err
		}

		re += n
		try++
		if n == nr {
			break
		}

	}

	if re != nr {
		return fmt.Errorf("short read: %d != %d", re, nr)
	}

	return nil
}

func write(r io.Writer, src []byte, nr int) error {

	return nil
}

func Encrypt(dest io.Writer, src io.Reader, key *[32]byte) (written uint64, err error) {
	buf := make([]byte, MaxPackageSize-NonceLength-secretbox.Overhead)

	// write the header
	if e := writeHeader(dest); e != nil {
		return 0, e
	}

	// Already used nounces
	var nonces [][]byte
	var msg []byte
	var nr int
	//var err error

	packages := 0

	for {

		if nr, err = src.Read(buf); err != nil {
			break
		}

		if nr > 0 {

			if msg, err = EncryptMessage(buf[0:nr], key, &nonces); err != nil {
				break
			}

			// Write package length
			if err = binary.Write(dest, binary.LittleEndian, uint16(len(msg))); err != nil {
				break
			}

			written += 2 // package length int16

			nw, ew := dest.Write(msg)

			written += uint64(nw)

			if ew != nil {
				err = ew
				break
			}

			if nw != len(msg) {
				err = io.ErrShortWrite
				break
			}

			packages += 1
		}

	}

	if err == io.EOF {
		err = nil
	}

	log.Printf("packages: %d\n", packages)

	return written, err

}

func Decrypt(dest io.Writer, src io.Reader, key *[32]byte) (err error) {

	var nr int

	var header [HeaderLength]byte

	// Read a validate header
	if nr, err = src.Read(header[:]); err != nil {
		return err
	} else if nr != HeaderLength {
		return errors.New("header length")
	}

	if string(header[:]) != "fnc" {
		return errors.New("file format")
	}

	// Read file
	var pkgSize uint16
	var msg []byte
	buf := make([]byte, MaxUint16)
	for {

		if err = binary.Read(src, binary.LittleEndian, &pkgSize); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("read packagesize: %s", err)
		}

		segSize := int(pkgSize) // int(pkgSize) + NonceLength + secretbox.Overhead

		/*if nr, err = src.Read(buf[0:segSize]); err != nil {
			return err
		} else if nr != segSize {
			return errors.New("short read")
		}*/

		if err = read(src, buf[0:segSize], segSize); err != nil && err != io.EOF {
			return err
		}

		if msg, err = DecryptMessage(buf[0:segSize], key); err != nil {
			return err
		}

		if nr, err = dest.Write(msg); err != nil {
			return err
		} else if nr != len(msg) {
			return fmt.Errorf("short write: %d != %d", nr, pkgSize)
		}

	}

	if err == io.EOF {
		err = nil
	}

	return err
}
