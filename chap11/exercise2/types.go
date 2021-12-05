package main

type pkgRegisterResponse struct {
	ID string `json:"id"`
}

type pkgQueryParams struct {
	name    string
	version string
	ownerId int
}

type pkgRow struct {
	OwnerId       int    `json:"owner_id"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	ObjectStoreId string `json:"object_store_id"`
	Created       string `json:"created"`
}

type pkgQueryResponse struct {
	Packages []pkgRow `json:"packages"`
}
