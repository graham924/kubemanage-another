package kube

import (
	"context"
	"github.com/pkg/errors"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NameSpace 全局变量，操作kubernetes环境使用
var NameSpace namespace

// namespace 自定义操作ns的类
type namespace struct{}

// NameSpaceResp 获取ns响应体数据结构
type NameSpaceResp struct {
	// Total 数量
	Total int `json:"total"`
	// Items coreV1.Namespace是k8s原生的ns结构，包括metadata、spec、status
	Items []coreV1.Namespace `json:"items"`
}

func (n *namespace) toCells(nodes []coreV1.Namespace) []DataCell {
	cells := make([]DataCell, len(nodes))
	for i := range nodes {
		cells[i] = namespaceCell(nodes[i])
	}
	return cells
}

func (n *namespace) FromCells(cells []DataCell) []coreV1.Namespace {
	nodes := make([]coreV1.Namespace, len(cells))
	for i := range cells {
		nodes[i] = coreV1.Namespace(cells[i].(namespaceCell))
	}
	return nodes
}

// GetNameSpaces 从k8s中获取ns列表
func (n *namespace) GetNameSpaces(filterName string, limit, page int) (nodesResp *NameSpaceResp, err error) {
	// 拿到全局的K8s客户端，获取内置的clientSet，就可以获取环境中的ns列表，返回的是k8s原生的nslist结构
	NamespaceList, err := K8sCli.ClientSet.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, errors.New("获取Namespace列表失败")
	}

	// 实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: n.toCells(NamespaceList.Items),
		DataSelect: &DataSelectQuery{
			Filter:     &FilterQuery{Name: filterName},
			Paginatite: &PaginateQuery{limit, page},
		},
	}
	// 先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	// 排序、分页
	data := filtered.Sort().Paginate()
	// 将 dataCell列表 转换为 coreV1.Namespace列表
	namespaces := n.FromCells(data.GenericDataList)
	// 返回的是封装好的响应体
	return &NameSpaceResp{
		total,
		namespaces,
	}, nil
}

// GetNameSpacesDetail 获取Node详情
func (n *namespace) GetNameSpacesDetail(Name string) (*coreV1.Namespace, error) {
	namespacesRes, err := K8sCli.ClientSet.CoreV1().Namespaces().Get(context.TODO(), Name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return namespacesRes, nil
}

func (n *namespace) CreateNameSpace(name string) error {
	ns := &coreV1.Namespace{
		TypeMeta: metaV1.TypeMeta{},
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
		},
		Spec:   coreV1.NamespaceSpec{},
		Status: coreV1.NamespaceStatus{},
	}
	if _, err := K8sCli.ClientSet.CoreV1().Namespaces().Create(context.TODO(), ns, metaV1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

func (n *namespace) DeleteNameSpace(name string) error {
	return K8sCli.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), name, metaV1.DeleteOptions{})
}
