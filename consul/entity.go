package consul

type host struct {
	Node struct {
		Node string
	}
	Service struct {
		Meta        map[string]string
		CreateIndex int
		ModifyIndex int
	}
}

type rawHosts []host
type hostsMap map[string]host
