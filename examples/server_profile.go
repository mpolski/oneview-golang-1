package main

import (
	"fmt"
	"github.com/HewlettPackard/oneview-golang/ov"
	"github.com/HewlettPackard/oneview-golang/utils"
	"os"
)

func main() {
	var (
		clientOV    *ov.OVClient
		sp_name     = "test"
		sp_sn       = "VCGRE1S007"
		new_sp_name = "Renamed Server Profile"
	)
	ovc := clientOV.NewOVClient(
		os.Getenv("ONEVIEW_OV_USER"),
		os.Getenv("ONEVIEW_OV_PASSWORD"),
		os.Getenv("ONEVIEW_OV_DOMAIN"),
		os.Getenv("ONEVIEW_OV_ENDPOINT"),
		false,
		800,
		"*")

	initialScopeUris := new([]utils.Nstring)
	*initialScopeUris = append(*initialScopeUris, utils.NewNstring("/rest/scopes/74877630-9a22-4061-9db4-d12b6c4cfee0"))

	server_profile_create_map := ov.ServerProfile{
		Type:              "ServerProfileV9",
		Name:              sp_name,
		ServerHardwareURI: "/rest/server-hardware/36343537-3338-4E43-3735-3532304D315A",
		InitialScopeUris:  *initialScopeUris,
	}

	err := ovc.SubmitNewProfile(server_profile_create_map)
	if err != nil {
		fmt.Println("Server Profile Create Failed: ", err)
	} else {
		fmt.Println("#----------------Server Profile Created---------------#")
	}

	sort := ""

	spt, err := ovc.GetProfileTemplateByName("go_test")
	if err != nil {
		fmt.Println("Server Profile Template Retrieval By Name Failed: ", err)
	} else {
		serverName, err := ovc.GetServerHardwareByName("DL980p Gen8 1 (new name)")
		if err != nil {
			fmt.Println("Failed to fetch server hardware name: ", err)
		} else {
			err = ovc.CreateProfileFromTemplate(sp_name, spt, serverName)
			if err != nil {
				fmt.Println("Server Profile Create Failed: ", err)
			} else {
				fmt.Println("#----------------Server Profile Created---------------#")
			}
		}
	}

	sp_list, err := ovc.GetProfiles("", "", "", sort, "")
	if err != nil {
		fmt.Println("Server Profile Retrieval Failed: ", err)
	} else {
		fmt.Println("#----------------Server Profile List---------------#")

		for i := 0; i < len(sp_list.Members); i++ {
			fmt.Println(sp_list.Members[i].Name)
		}
	}

	sp1, err := ovc.GetProfileByName(sp_name)
	if err != nil {
		fmt.Println("Server Profile Retrieval By Name Failed: ", err)
	} else {
		fmt.Println("#----------------Server Profile by Name---------------#")
		fmt.Println(sp1.Name)
	}

	sp2, err := ovc.GetProfileBySN(sp_sn)
	if err != nil {
		fmt.Println("Server Profile Retrieval By Serial Number Failed: ", err)
	} else {
		fmt.Println("#----------------Server Profile by Serial Number---------------#")
		fmt.Println(sp2.Name)
	}

	sp, err := ovc.GetProfileByURI(sp2.URI)
	if err != nil {
		fmt.Println("Server Profile Retrieval By URI Failed: ", err)
	} else {
		fmt.Println("#----------------Server Profile by URI---------------#")
		fmt.Println(sp.Name)
	}

	sp_update_clone := ov.ServerProfile{
		Name:                  new_sp_name,
		URI:                   sp1.URI,
		Type:                  sp1.Type,
		ETAG:                  sp1.ETAG,
		Affinity:              sp1.Affinity,
		ServerHardwareTypeURI: sp1.ServerHardwareTypeURI,
		EnclosureGroupURI:     sp1.EnclosureGroupURI,
	}

	err = ovc.UpdateServerProfile(sp_update_clone)
	if err != nil {
		fmt.Println("Server Profile Create Failed: ", err)
	} else {
		fmt.Println("#----------------Server Profile Created---------------#")
	}

	sp_list, err = ovc.GetProfiles("", "", "", sort, "")
	if err != nil {
		fmt.Println("Server Profile Retrieval Failed: ", err)
	} else {
		fmt.Println("#----------------Server Profile List---------------#")

		for i := 0; i < len(sp_list.Members); i++ {
			fmt.Println(sp_list.Members[i].Name)
		}
	}

	task, err := ovc.SubmitDeleteProfile(sp1)
	if err != nil {
		fmt.Println("Server Profile Delete Request Failed: ", err)
	} else {
		fmt.Println("#----------------Server Profile Delete---------------#")
		fmt.Println("Task URI: ", task.URI)
	}

	err = ovc.DeleteProfile(new_sp_name)
	if err != nil {
		fmt.Println("Server Profile Delete Failed: ", err)
	} else {
		fmt.Println("#----------------Server Profile Deleted---------------#")
	}

}
