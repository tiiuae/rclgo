package rclgo

/*
#include <rcl_action/graph.h>
#include <rcl/graph.h>
*/
import "C"

import (
	"unsafe"
)

// GetTopicNamesAndTypes returns a map of all known topic names to corresponding
// topic types. Note that multiple types may be associated with a single topic.
//
// If demangle is true, topic names will be in the format used by the underlying
// middleware.
func (n *Node) GetTopicNamesAndTypes(demangle bool) (map[string][]string, error) {
	return n.getNamesAndTypes("", "", func(node, namespace *C.char, namesAndTypes *C.rmw_names_and_types_t) C.int {
		return C.rcl_get_topic_names_and_types(
			n.rcl_node_t,
			n.context.rcl_allocator_t,
			!C.bool(demangle),
			namesAndTypes,
		)
	})
}

func (n *Node) GetNodeNames() (names, namespaces []string, err error) {
	rcnames := C.rcutils_get_zero_initialized_string_array()
	defer C.rcutils_string_array_fini(&rcnames)
	rcnamespaces := C.rcutils_get_zero_initialized_string_array()
	defer C.rcutils_string_array_fini(&rcnamespaces)
	rc := C.rcl_get_node_names(
		n.rcl_node_t,
		*n.context.rcl_allocator_t,
		&rcnames,
		&rcnamespaces,
	)
	if rc != C.RCL_RET_OK {
		return nil, nil, errorsCastC(rc, "failed to get node names")
	}
	cnames := unsafe.Slice(rcnames.data, rcnames.size)
	names = make([]string, rcnames.size)
	cnamespaces := unsafe.Slice(rcnamespaces.data, rcnamespaces.size)
	namespaces = make([]string, rcnamespaces.size)
	for i := range names {
		names[i] = C.GoString(cnames[i])
		namespaces[i] = C.GoString(cnamespaces[i])
	}
	return names, namespaces, nil
}

func (n *Node) GetPublisherNamesAndTypesByNode(demangle bool, node, namespace string) (map[string][]string, error) {
	return n.getNamesAndTypes(node, namespace, func(node, namespace *C.char, namesAndTypes *C.rmw_names_and_types_t) C.int {
		return C.rcl_get_publisher_names_and_types_by_node(
			n.rcl_node_t,
			n.context.rcl_allocator_t,
			!C.bool(demangle),
			node,
			namespace,
			namesAndTypes,
		)
	})
}

func (n *Node) GetSubscriberNamesAndTypesByNode(demangle bool, node, namespace string) (map[string][]string, error) {
	return n.getNamesAndTypes(node, namespace, func(node, namespace *C.char, namesAndTypes *C.rmw_names_and_types_t) C.int {
		return C.rcl_get_subscriber_names_and_types_by_node(
			n.rcl_node_t,
			n.context.rcl_allocator_t,
			!C.bool(demangle),
			node,
			namespace,
			namesAndTypes,
		)
	})
}

func (n *Node) GetServiceNamesAndTypesByNode(node, namespace string) (map[string][]string, error) {
	return n.getNamesAndTypes(node, namespace, func(node, namespace *C.char, namesAndTypes *C.rmw_names_and_types_t) C.int {
		return C.rcl_get_service_names_and_types_by_node(
			n.rcl_node_t,
			n.context.rcl_allocator_t,
			node,
			namespace,
			namesAndTypes,
		)
	})
}

func (n *Node) GetClientNamesAndTypesByNode(node, namespace string) (map[string][]string, error) {
	return n.getNamesAndTypes(node, namespace, func(node, namespace *C.char, namesAndTypes *C.rmw_names_and_types_t) C.int {
		return C.rcl_get_client_names_and_types_by_node(
			n.rcl_node_t,
			n.context.rcl_allocator_t,
			node,
			namespace,
			namesAndTypes,
		)
	})
}

func (n *Node) GetActionServerNamesAndTypesByNode(node, namespace string) (map[string][]string, error) {
	return n.getNamesAndTypes(node, namespace, func(node, namespace *C.char, namesAndTypes *C.rmw_names_and_types_t) C.int {
		return C.rcl_action_get_server_names_and_types_by_node(
			n.rcl_node_t,
			n.context.rcl_allocator_t,
			node,
			namespace,
			namesAndTypes,
		)
	})
}

func (n *Node) GetActionClientNamesAndTypesByNode(node, namespace string) (map[string][]string, error) {
	return n.getNamesAndTypes(node, namespace, func(node, namespace *C.char, namesAndTypes *C.rmw_names_and_types_t) C.int {
		return C.rcl_action_get_client_names_and_types_by_node(
			n.rcl_node_t,
			n.context.rcl_allocator_t,
			node,
			namespace,
			namesAndTypes,
		)
	})
}

func (n *Node) getNamesAndTypes(
	node, namespace string,
	get func(
		node, namespace *C.char,
		namesAndTypes *C.rmw_names_and_types_t,
	) C.int,
) (map[string][]string, error) {
	cnode := C.CString(node)
	defer C.free(unsafe.Pointer(cnode))
	cnamespace := C.CString(namespace)
	defer C.free(unsafe.Pointer(cnamespace))
	cnamesAndTypes := C.rcl_get_zero_initialized_names_and_types()
	defer C.rcl_names_and_types_fini(&cnamesAndTypes)
	rc := get(cnode, cnamespace, &cnamesAndTypes)
	if rc != C.RCL_RET_OK {
		return nil, errorsCastC(rc, "failed to get topic names and types")
	}
	names := unsafe.Slice(cnamesAndTypes.names.data, cnamesAndTypes.names.size)
	types := unsafe.Slice(cnamesAndTypes.types, len(names))
	namesAndTypes := make(map[string][]string, len(names))
	for i, name := range names {
		name := C.GoString(name)
		typesForName := unsafe.Slice(types[i].data, types[i].size)
		resultTypes := make([]string, len(typesForName))
		for j, typ := range typesForName {
			resultTypes[j] = C.GoString(typ)
		}
		namesAndTypes[name] = resultTypes
	}
	return namesAndTypes, nil
}

const GIDSize = 24

type GID [GIDSize]byte

type EndpointType int

const (
	EndpointInvalid EndpointType = iota
	EndpointPublisher
	EndpointSubscription
)

type TopicEndpointInfo struct {
	NodeName      string
	NodeNamespace string
	TopicType     string
	EndpointType  EndpointType
	EndpointGID   GID
	QosProfile    QosProfile
}

func (n *Node) GetPublishersInfoByTopic(topic string, mangle bool) ([]TopicEndpointInfo, error) {
	return n.getInfoByTopic("publishers", topic, mangle, func(topic *C.char, noMangle C.bool, infoArray *C.rmw_topic_endpoint_info_array_t) C.int {
		return C.rcl_get_publishers_info_by_topic(n.rcl_node_t, n.context.rcl_allocator_t, topic, noMangle, infoArray)
	})
}

func (n *Node) GetSubscriptionsInfoByTopic(topic string, mangle bool) ([]TopicEndpointInfo, error) {
	return n.getInfoByTopic("subscriptions", topic, mangle, func(topic *C.char, noMangle C.bool, infoArray *C.rmw_topic_endpoint_info_array_t) C.int {
		return C.rcl_get_subscriptions_info_by_topic(n.rcl_node_t, n.context.rcl_allocator_t, topic, noMangle, infoArray)
	})
}

func (n *Node) getInfoByTopic(kind, topic string, mangle bool, get func(
	topic *C.char,
	noMangle C.bool,
	infoArray *C.rmw_topic_endpoint_info_array_t,
) C.int) ([]TopicEndpointInfo, error) {
	ctopic := C.CString(topic)
	defer C.free(unsafe.Pointer(ctopic))
	infoArray := C.rmw_get_zero_initialized_topic_endpoint_info_array()
	defer C.rmw_topic_endpoint_info_array_fini(&infoArray, n.context.rcl_allocator_t)
	rc := get(ctopic, !C.bool(mangle), &infoArray)
	if rc != C.RCL_RET_OK {
		return nil, errorsCastC(rc, "failed to get "+kind+" info by topic")
	}
	infoSlice := unsafe.Slice(infoArray.info_array, infoArray.size)
	infos := make([]TopicEndpointInfo, len(infoSlice))
	for i, info := range infoSlice {
		infos[i] = TopicEndpointInfo{
			NodeName:      C.GoString(info.node_name),
			NodeNamespace: C.GoString(info.node_namespace),
			TopicType:     C.GoString(info.topic_type),
			EndpointType:  EndpointType(info.endpoint_type),
		}
		for j := range info.endpoint_gid {
			infos[i].EndpointGID[j] = byte(info.endpoint_gid[j])
		}
		infos[i].QosProfile.fromCStruct(&info.qos_profile)
	}
	return infos, nil
}
