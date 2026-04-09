package main

func AmazonEC2MachineID() (uint16, error) {
	ip, err := amazonEC2PrivateIPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}
