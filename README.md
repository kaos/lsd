Livestatusd
===========

A livestatus daemon supporting any backend monitoring system using plugins.

This is a Go implementation of
[livestatus](https://mathias-kettner.de/checkmk_livestatus.html) by
Mathias Kettner, but with support for more monitoring systems besides
Nagios.


Compatibility
-------------

The goal is to stay 100% compatible with the original livestatus
implementation. Some fields may however be missing, or behave slightly
different depending on the backend system sourcing the data.

Notable differences compared to the original livestatus that livestatusd supports:

  * multiple grouping columns in a stats request.

  * both unix and tcp socket support (a single instance may listen on
    both types at the same time).

  * introduced `info` table to retrive info about the livestatusd
    instance.



Current Status
--------------

This is work in progress, and not all features have yet been implemented.


Supported Backend Systems
-------------------------

Livestatusd supports collecting monitoring data from the following systems:

  * Sensu





Dev Notes
=========

Sample output
-------------

From the original livestatus, running against nagios:

    [root@host ~]# echo -e "GET hosts\nLimit: 1\nOutputFormat: json\nColumns: name groups plugin_output services_with_info" | unixcat /var/nagios/rw/live | jq .
    [
      [
        "foo.example.se",
        [
          "infra",
          "dc1_monitoring_hostgroup"
        ],
        "OK - foo.example.se responds to ICMP. Packet 1, rta 4.508ms",
        [
          [
            "Its alive",
            0,
            1,
            "PING OK - Packet loss = 0%, RTA = 1.32 ms"
          ]
        ]
      ]
    ]
    [root@host ~]# echo -e "GET hosts\nLimit: 1\nOutputFormat: csv\nColumns: name groups plugin_output services_with_info" | unixcat /var/nagios/rw/live
    foo.example.se;infra,dc1_monitoring_hostgroup;OK - foo.example.se responds to ICMP. Packet 1, rta 4.508ms;Its alive|0|1|PING OK - Packet loss = 0%, RTA = 1.32 ms


And python output:

    [root@host ~]# echo -e "GET hosts\nLimit: 1\nOutputFormat: python\nColumns: name groups plugin_output services_with_info" | unixcat /var/nagios/rw/live
    [[u"foo.example.se",[u"infra",u"dc1_monitoring_hostgroup"],u"OK - foo.example.se responds to ICMP. Packet 1, rta 4.508ms",[[u"Its alive",0,1,u"PING OK - Packet loss = 0%, RTA = 1.32 ms"]]]]



Error handling
--------------

Simple, any invalid input yields no output.



Fields
------

## Hosts

    [root@host ~]# echo -e "GET hosts\nLimit: 0\nOutputFormat: csv" | unixcat /var/nagios/rw/live
    accept_passive_checks;acknowledged;acknowledgement_type;action_url;action_url_expanded;active_checks_enabled;address;alias;check_command;check_flapping_recovery_notification;check_freshness;check_interval;check_options;check_period;check_source;check_type;checks_enabled;childs;comments;comments_with_info;contact_groups;contacts;current_attempt;current_notification_number;custom_variable_names;custom_variable_values;custom_variables;display_name;downtimes;downtimes_with_info;event_handler;event_handler_enabled;execution_time;filename;first_notification_delay;flap_detection_enabled;groups;hard_state;has_been_checked;high_flap_threshold;hourly_value;icon_image;icon_image_alt;icon_image_expanded;id;in_check_period;in_notification_period;initial_state;is_executing;is_flapping;last_check;last_hard_state;last_hard_state_change;last_notification;last_state;last_state_change;last_time_down;last_time_unreachable;last_time_up;latency;long_plugin_output;low_flap_threshold;max_check_attempts;modified_attributes;modified_attributes_list;name;next_check;next_notification;no_more_notifications;notes;notes_expanded;notes_url;notes_url_expanded;notification_interval;notification_period;notifications_enabled;num_services;num_services_crit;num_services_hard_crit;num_services_hard_ok;num_services_hard_unknown;num_services_hard_warn;num_services_ok;num_services_pending;num_services_unknown;num_services_warn;obsess;parents;pending_flex_downtime;percent_state_change;perf_data;plugin_output;pnpgraph_present;process_performance_data;retry_interval;scheduled_downtime_depth;services;services_with_info;services_with_state;should_be_scheduled;state;state_type;statusmap_image;total_services;worst_service_hard_state;worst_service_state;x_3d;y_3d;z_3d

## Services

    [root@host ~]# echo -e "GET services\nLimit: 0\nOutputFormat: csv" | unixcat /var/nagios/rw/live
    accept_passive_checks;acknowledged;acknowledgement_type;action_url;action_url_expanded;active_checks_enabled;check_command;check_freshness;check_interval;check_options;check_period;check_source;check_type;checks_enabled;comments;comments_with_info;contact_groups;contacts;current_attempt;current_notification_number;custom_variable_names;custom_variable_values;custom_variables;description;display_name;downtimes;downtimes_with_info;event_handler;event_handler_enabled;execution_time;first_notification_delay;flap_detection_enabled;groups;has_been_checked;high_flap_threshold;host_accept_passive_checks;host_acknowledged;host_acknowledgement_type;host_action_url;host_action_url_expanded;host_active_checks_enabled;host_address;host_alias;host_check_command;host_check_flapping_recovery_notification;host_check_freshness;host_check_interval;host_check_options;host_check_period;host_check_source;host_check_type;host_checks_enabled;host_childs;host_comments;host_comments_with_info;host_contact_groups;host_contacts;host_current_attempt;host_current_notification_number;host_custom_variable_names;host_custom_variable_values;host_custom_variables;host_display_name;host_downtimes;host_downtimes_with_info;host_event_handler;host_event_handler_enabled;host_execution_time;host_filename;host_first_notification_delay;host_flap_detection_enabled;host_groups;host_hard_state;host_has_been_checked;host_high_flap_threshold;host_hourly_value;host_icon_image;host_icon_image_alt;host_icon_image_expanded;host_id;host_in_check_period;host_in_notification_period;host_initial_state;host_is_executing;host_is_flapping;host_last_check;host_last_hard_state;host_last_hard_state_change;host_last_notification;host_last_state;host_last_state_change;host_last_time_down;host_last_time_unreachable;host_last_time_up;host_latency;host_long_plugin_output;host_low_flap_threshold;host_max_check_attempts;host_modified_attributes;host_modified_attributes_list;host_name;host_next_check;host_next_notification;host_no_more_notifications;host_notes;host_notes_expanded;host_notes_url;host_notes_url_expanded;host_notification_interval;host_notification_period;host_notifications_enabled;host_num_services;host_num_services_crit;host_num_services_hard_crit;host_num_services_hard_ok;host_num_services_hard_unknown;host_num_services_hard_warn;host_num_services_ok;host_num_services_pending;host_num_services_unknown;host_num_services_warn;host_obsess;host_parents;host_pending_flex_downtime;host_percent_state_change;host_perf_data;host_plugin_output;host_pnpgraph_present;host_process_performance_data;host_retry_interval;host_scheduled_downtime_depth;host_services;host_services_with_info;host_services_with_state;host_should_be_scheduled;host_state;host_state_type;host_statusmap_image;host_total_services;host_worst_service_hard_state;host_worst_service_state;host_x_3d;host_y_3d;host_z_3d;hourly_value;icon_image;icon_image_alt;icon_image_expanded;id;in_check_period;in_notification_period;initial_state;is_executing;is_flapping;last_check;last_hard_state;last_hard_state_change;last_notification;last_state;last_state_change;last_time_critical;last_time_ok;last_time_unknown;last_time_warning;latency;long_plugin_output;low_flap_threshold;max_check_attempts;modified_attributes;modified_attributes_list;next_check;next_notification;no_more_notifications;notes;notes_expanded;notes_url;notes_url_expanded;notification_interval;notification_period;notifications_enabled;obsess;percent_state_change;perf_data;plugin_output;pnpgraph_present;process_performance_data;retry_interval;scheduled_downtime_depth;should_be_scheduled;state;state_type
