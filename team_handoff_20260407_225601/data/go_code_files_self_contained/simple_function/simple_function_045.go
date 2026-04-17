package main

func Repository() (string, error) {
	url, err := OriginURL()
	if err != nil {
		return "", err
	}

	repo := retrieveRepoName(url)
	return repo, nil
}
