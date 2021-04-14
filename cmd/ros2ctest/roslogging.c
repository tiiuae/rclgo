#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include <rcl/arguments.h>
#include <rcl/init.h>
#include <rcl/logging.h>
#include <rcl/node.h>
#include <rcl/subscription.h>
#include <rcutils/error_handling.h>

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <std_msgs/msg/color_rgba.h>

void init_logging(int argc, const char** argv) {
    rcl_ret_t rc;
    const rcl_allocator_t allocator = rcl_get_default_allocator();

    rcl_arguments_t rcl_arguments = rcl_get_zero_initialized_arguments();
    rc = rcl_parse_arguments(argc, argv, allocator, &rcl_arguments);
    if (rc != RCL_RET_OK) {
        printf("rcl_parse_arguments error '%d' '%s'\n", rc, rcutils_get_error_string().str);
        rcutils_reset_error();
    }

    rc = rcl_logging_configure(&rcl_arguments, &allocator);
    if (rc != RCL_RET_OK) {
        printf("rcl_logging_configure error '%d' '%s'\n", rc, rcutils_get_error_string().str);
        rcutils_reset_error();
    }
}

int main(int argc, const char** argv) {
    rcl_ret_t rc;
    const rcl_allocator_t allocator = rcl_get_default_allocator();

    // Overwrite argv with ROS2 logging settings
    printf("argc %d argv%s\n", argc, *argv);
    const char * newArgv [] = {"--ros-args", "--log-level", "DEBUG"};
    argv = newArgv;
    argc = 3;
    printf("argc %d argv%s\n", argc, *argv);

    init_logging(argc, argv);

    rcl_context_t context = rcl_get_zero_initialized_context();
    rcl_init_options_t options = rcl_get_zero_initialized_init_options();
    rc = rcl_init_options_init(&options, allocator);

    rc = rcl_init(argc, argv, &options, &context);
    if (rc != RCL_RET_OK) {
        printf("rcl_init error '%d' '%s'\n", rc, rcutils_get_error_string().str);
        rcutils_reset_error();
    }

    rcl_node_options_t node_options = rcl_node_get_default_options();
    rcl_node_t node = rcl_get_zero_initialized_node();
    rc = rcl_node_init(&node, "NODE_NAME12", "/", &context, &node_options);
    if (rc != RCL_RET_OK) {
        printf("rcl_node_init error '%d' '%s'\n", rc, rcutils_get_error_string().str);
        rcutils_reset_error();
    }

    rcl_subscription_t subscription = rcl_get_zero_initialized_subscription();
    rcl_subscription_options_t subscription_options = rcl_subscription_get_default_options();
    rc = rcl_subscription_init(&subscription, &node, rosidl_typesupport_c__get_message_type_support_handle__std_msgs__msg__ColorRGBA(), "/rossina/rusina", &subscription_options);
    if (rc != RCL_RET_OK) {
        printf("rcl_subscription_init error '%d' '%s'\n", rc, rcutils_get_error_string().str);
        rcutils_reset_error();
    }

	rmw_message_info_t rmw_message_info = rmw_get_zero_initialized_message_info();
	std_msgs__msg__ColorRGBA* ros2_msg_receive_buffer = std_msgs__msg__ColorRGBA__create();

    sleep(2);

    rc = rcl_take(&subscription, ros2_msg_receive_buffer, &rmw_message_info, NULL);
    if (rc != RCL_RET_OK) {
        printf("rcl_take error '%d' '%s'\n", rc, rcutils_get_error_string().str);
        rcutils_reset_error();
    }

    printf("ColorRGBA: %f %f %f %f\n", ros2_msg_receive_buffer->r, ros2_msg_receive_buffer->b, ros2_msg_receive_buffer->g, ros2_msg_receive_buffer->a);

    std_msgs__msg__ColorRGBA__fini(ros2_msg_receive_buffer);

    rc = rcl_subscription_fini(&subscription, &node);
    if (rc != RCL_RET_OK) {
        printf("rcl_subscription_fini error '%d' '%s'\n", rc, rcutils_get_error_string().str);
        rcutils_reset_error();
    }
    rc = rcl_node_fini(&node);
    if (rc != RCL_RET_OK) {
        printf("rcl_node_fini error '%d' '%s'\n", rc, rcutils_get_error_string().str);
        rcutils_reset_error();
    }
    rc = rcl_shutdown(&context);
    if (rc != RCL_RET_OK) {
        printf("rcl_shutdown error '%d' '%s'\n", rc, rcutils_get_error_string().str);
        rcutils_reset_error();
    }
    rc = rcl_logging_fini();
    if (rc != RCL_RET_OK) {
        printf("rcl_logging_fini error '%d' '%s'\n", rc, rcutils_get_error_string().str);
        rcutils_reset_error();
    }

    return 0;
}
