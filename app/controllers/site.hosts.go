package controllers

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func (s *Site) ReadHostsFile(hostsFile string) (models.HostsMap, error) {
	fi, err := os.Open(hostsFile)

	if err != nil {
		log.Println(err)
		return models.HostsMap{}, err
	}

	defer fi.Close()

	return s.ParseHostsContent(fi)
}

func (s *Site) ParseHostsContent(content io.Reader) (models.HostsMap, error) {
	hostsMap := models.NewHostsMap()
	var (
		// Make a read buffer
		r         = bufio.NewReader(content)
		readState = models.READ_NOT_STARTED
	)

	var str string
	for {
		ln, prefixed, err := r.ReadLine()

		// If error wasn't nil return error.
		if err != nil && err != io.EOF {
			log.Println(err)
			return hostsMap, err
		}

		// If no bytes where read return from loop.
		if err == io.EOF {
			break
		}
		// Append bytes (as string) to string.
		str += string(ln)

		// If prefixed continue, get rest of the bytes from line..
		if prefixed {
			continue
		}

		// Trim spaces from string.
		str = strings.TrimSpace(str)

		// Get application hosts read start.
		if str == models.HOSTS_START_READ_STR {
			log.Println("Application hosts read started.")
			readState = models.READ_STARTED
			str = ""
			continue
		}

		// Get application hosts read end and brake.
		if str == models.HOSTS_END_READ_STR {
			log.Println("Application hosts read ended.")
			readState = models.READ_ENDED
			break
		}

		if readState == models.READ_STARTED {
			hd := models.NewHostDomains()
			err := hd.Parse(str)
			if err != nil {
				log.Println(err)
			}
			hostsMap.AddHostDomains(hd)
		}

		str = ""
	}

	if readState != models.READ_ENDED {
		msg := fmt.Sprintf("No application hosts found on hosts file.\n")
		return hostsMap, errors.New(msg)
	}

	return hostsMap, nil
}

func (s *Site) WriteNewHosts(hostsFile string, hostsMap *models.HostsMap) error {
	fi, err := os.Open(hostsFile)
	if err != nil {
		log.Println(err)
		return err
	}
	r := bufio.NewReader(fi)

	var (
		// Create temporary slice of bytes.
		temp      []byte
		readState = models.READ_NOT_STARTED
	)
	for {
		// Prefixed indicates if the current line was only partially read.
		// If is prefixed we won't later add new line byte to slice.
		ln, prefixed, err := r.ReadLine()
		str := string(ln)

		// If error wasn't nil return error.
		if err != nil && err != io.EOF {
			log.Println(err)
			return err
		}

		// If no bytes where read return from loop.
		if err == io.EOF {
			break
		}

		if str == models.HOSTS_START_READ_STR {
			// Set read state started.
			readState = models.READ_STARTED
			// Write hosts map as bytes.
			temp = append(temp, hostsMap.Bytes(models.SPACE_BYTE, models.NEW_LINE_BYTE)...)
		}

		if str == models.HOSTS_END_READ_STR {
			readState = models.READ_ENDED
			// Empty the string so we don't add snippet end line twice.
			ln = ln[:0]
			str = ""
		}

		// If read started is indicated, don't write contents.
		// Hosts in snippet area are rebuild.
		if readState != models.READ_STARTED {
			if len(ln) > 0 {
				temp = append(temp, ln...)
				// If prefixed we don't add new line byte to slice.
				if prefixed {
					continue
				}

				temp = append(temp, models.NEW_LINE_BYTE)
			}
		}
	}
	ok, err := s.HostsContentAssertsTrue(bytes.NewReader(temp))
	if !ok {
		return err
	}

	_, err = utils.CreateBackupFile(hostsFile)
	err = ioutil.WriteFile(hostsFile, temp, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Site) HostsContentAssertsTrue(content io.Reader) (bool, error) {
	_, err := s.ParseHostsContent(content)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Site) AddToHosts(templ *models.InstallTemplate, domains *models.SiteDomains) error {
	hostsFile := "/etc/hosts"

	hostsMap, err := s.ReadHostsFile(hostsFile)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, domain := range domains.Domains {
		hostsMap.AddDomain(domain)
	}

	err = s.WriteNewHosts(hostsFile, &hostsMap)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
