package main

import (
	"fmt"
	"log"
	"github.com/trivago/grok"
)
 
func main() {
	//  fmt.Println("# Default Capture :")
	//  g, _ := grok.New(grok.Config{})
	//  values, _ := g.ParseString("%{COMMONAPACHELOG}", `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`)
	//  for k, v := range values {
	// 	 fmt.Printf("%+15s: %s\n", k, v)
	//  }
 
	//  fmt.Println("\n# Named Capture :")
	//  g, _ = grok.New(grok.Config{NamedCapturesOnly: true})
	//  values, _ = g.ParseString("%{COMMONAPACHELOG}", `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`)
	//  for k, v := range values {
	// 	 fmt.Printf("%+15s: %s\n", k, v)
	//  }
 
	//  fmt.Println("\n# Add custom patterns :")
	//  // We add 3 patterns to our Grok instance, to structure an IRC message

	pBasicFields := map[string]string{
		"UUID": `[a-z0-9]{8}-([a-z0-9]{4}-){3}[a-z0-9]{12}`,    	// 6245c77d-5017-4657-b35b-7ab1d247112b
		"REQID": `req-%{UUID}`,										// req-8cadad28-8315-45ca-818c-6a229dfb73e1
		"DATETIME": `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}`,	// 2020-09-27 19:22:54.486
		"MD5": `[0-9a-z]{32}`, 										// 62c38230485b4794a8eedece5dac9192
		"JSON": `\{.*\}`,											// {u'bandwidth_limit_rule': {u'max_kbps': 102400, u'direction': u'egress', u'max_burst_kbps': 102400}}
	}

	pCommon := map[string]string{
		// 2020-09-27 19:22:54.485 68316 DEBUG neutron.api.v2.base 
		// [req-8cadad28-8315-45ca-818c-6a229dfb73e1 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - default default] Request body: 
		// {u'bandwidth_limit_rule': {u'max_kbps': 102400, u'direction': u'egress', u'max_burst_kbps': 102400}} 
		// prepare_request_body /usr/lib/python2.7/site-packages/neutron/api/v2/base.py:713
		"neutron_api_v2_base": `%{DATETIME:neutron_api_time} .* neutron.api.v2.base \[%{REQID:req_id} .*\] ` +
							   `Request body: %{JSON:request_body} prepare_request_body .*$`,

	}

	pLoadBalancerCreate := map[string]string{
		// 05neu-core/server.log-1005:2020-10-05 10:20:17.251 117825 DEBUG f5lbaasdriver.v2.bigip.driver_v2 
		// [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - default default] 
		// f5lbaasdriver.v2.bigip.driver_v2.LoadBalancerManager method create called with arguments (<neutron_lib.context.Context object at 0x284cb250>, 
		// <neutron_lbaas.services.loadbalancer.data_models.LoadBalancer object at 0xdb44250>) {} 
		// wrapper /usr/lib/python2.7/site-packages/oslo_log/helpers.py:66
		"call_f5driver": 
			`%{DATETIME:call_f5driver_time} .* f5lbaasdriver.v2.bigip.driver_v2 \[%{REQID:req_id} .*\] ` +
			`f5lbaasdriver.v2.bigip.driver_v2.LoadBalancerManager method create called with .*$`,
		
		// 2020-10-05 10:20:21.924 117825 DEBUG f5lbaasdriver.v2.bigip.agent_scheduler 
		// [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - default default] 
		// Loadbalancer e2d277f7-eca2-46a4-bf2c-655856fd8733 is scheduled to lbaas agent dc55e196-319a-4c82-b262-344f45415131 schedule 
		// /usr/lib/python2.7/site-packages/f5lbaasdriver/v2/bigip/agent_scheduler.py:306
		// "agent_scheduled": 

		// 2020-10-05 10:20:27.176 117825 DEBUG f5lbaasdriver.v2.bigip.agent_rpc [req-92db71fb-8
		// 513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - default default]
		// f5lbaasdriver.v2.bigip.agent_rpc.LBaaSv2AgentRPC method create_loadbalancer called with arguments (<neutron_lib.
		// context.Context object at 0x284cb250>, {'availability_zone_hints': [], 'description': '', 'admin_state_up': True
		// , 'tenant_id': '62c38230485b4794a8eedece5dac9192', 'provisioning_status': 'PENDING_CREATE', 'listeners': [], 'vi
		// p_subnet_id': 'd79ef712-c1e3-4860-9343-d1702b9976aa', 'vip_address': '10.230.44.15', 'vip_port_id': '5bcbe2d7-99
		// 4f-40de-87ab-07aa632f0133', 'provider': None, 'pools': [], 'id': 'e2d277f7-eca2-46a4-bf2c-655856fd8733', 'operat
		// ing_status': 'OFFLINE', 'name': 'JL-B01-POD1-CORE-LB-7'}, {'subnets': ...
		// : 'd79ef712-c1e3-4860-9343-d1702b9976aa', 'vip_address': '10.230.44.15', 'vip_port_id': '5bcbe2d7-994f-40de-87ab
		// -07aa632f0133', 'provider': None, 'pools': [], 'id': 'e2d277f7-eca2-46a4-bf2c-655856fd8733', 'operating_status':
		// 'OFFLINE', 'name': 'JL-B01-POD1-CORE-LB-7'}}, u'POD1_CORE3') {} wrapper /usr/lib/python2.7/site-packages/oslo_l
		// og/helpers.py:66
		"rpc_f5agent": 
			`%{DATETIME:rpc_f5agent_time} .* f5lbaasdriver.v2.bigip.agent_rpc \[%{REQID:req_id} .*\] ` +
			`f5lbaasdriver.v2.bigip.agent_rpc.LBaaSv2AgentRPC method create_loadbalancer called with arguments ` +
			`.*? 'id': '%{UUID:object_id}'.*`,

		// 2020-10-05 10:19:16.315 295263 DEBUG f5_openstack_agent.lbaasv2.drivers.bigip.agent_manager 
		// [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - - -] 
		// f5_openstack_agent.lbaasv2.drivers.bigip.agent_manager.LbaasAgentManager method create_loadbalancer called with arguments
		// ...
		// 7'}} wrapper /usr/lib/python2.7/site-packages/oslo_log/helpers.py:66
		"call_f5agent": 
			`%{DATETIME:call_f5agent_time} .* f5_openstack_agent.lbaasv2.drivers.bigip.agent_manager \[%{REQID:req_id} .*\] ` +
			`f5_openstack_agent.lbaasv2.drivers.bigip.agent_manager.LbaasAgentManager method create_loadbalancer called with arguments .*`,

		// 2020-10-05 10:19:16.317 295263 DEBUG root [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 
		// 62c38230485b4794a8eedece5dac9192 - - -] get WITH uri: https://10.216.177.8:443/mgmt/tm/sys/folder/~CORE_62c38230485b4794a8eedece5dac9192 AND 
		// suffix:  AND kwargs: {} wrapper /usr/lib/python2.7/site-packages/icontrol/session.py:257
		"rest_call_bigip": 
			`%{DATETIME:call_bigip_time} .* \[%{REQID:req_id} .*\] get WITH uri: .*icontrol/session.py.*`,

		// 2020-10-05 10:19:18.411 295263 DEBUG f5_openstack_agent.lbaasv2.drivers.bigip.plugin_rpc 
		// [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - - -] 
		// f5_openstack_agent.lbaasv2.drivers.bigip.plugin_rpc.LBaaSv2PluginRPC method update_loadbalancer_status called with arguments 
		// (u'e2d277f7-eca2-46a4-bf2c-655856fd8733', 'ACTIVE', 'ONLINE', u'JL-B01-POD1-CORE-LB-7') {} wrapper 
		// /usr/lib/python2.7/site-packages/oslo_log/helpers.py:66
		"report_loadbalancer_status": 
			`%{DATETIME:call_bigip_time} .* f5_openstack_agent.lbaasv2.drivers.bigip.plugin_rpc \[%{REQID:req_id} .*\].* ` +
			`method update_loadbalancer_status called with arguments.*`,
	}

	// pMemberCreate := map[string]string{
	// 	"sdf": ``,
	// }

	pattern :=map[string]string {}

	for k, v := range pBasicFields {
		pattern[k] = v
	}
	for k, v := range pCommon {
		pattern[k] = v
	}

	for k, v := range pLoadBalancerCreate {
		pattern[k] = v
	}

	g, e := grok.New(grok.Config{
		NamedCapturesOnly: true,
		Patterns: pattern,
	})
	if e != nil {
		log.Panic(e)
	}

	for k, v := range tests() {
		fmt.Printf("------- %s --------\n", k)
		value, err := test_sg(k, v, g)
		debug(value, err)
	}
}

func debug(values map[string]string, e error) {
	if e != nil {
		log.Println(e.Error())
		return
	 }

	 for k, v := range values {
		 log.Printf("%+25s: %s\n", k, v)
	 }
	 log.Println()
}

func test_sg(k string, v string, g *grok.Grok) (map[string]string, error) {
	return g.ParseString(fmt.Sprintf("%%{%s}", k), v)
}

func tests() map[string]string {
	return map[string]string{
		"rpc_f5agent": 
			`2020-10-05 10:20:27.176 117825 DEBUG f5lbaasdriver.v2.bigip.agent_rpc [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - default default] f5lbaasdriver.v2.bigip.agent_rpc.LBaaSv2AgentRPC method create_loadbalancer called with arguments (<neutron_lib.context.Context object at 0x284cb250>, {'availability_zone_hints': [], 'description': '', 'admin_state_up': True, 'tenant_id': '62c38230485b4794a8eedece5dac9192', 'provisioning_status': 'PENDING_CREATE', 'listeners': [], 'vip_subnet_id': 'd79ef712-c1e3-4860-9343-d1702b9976aa', 'vip_address': '10.230.44.15', 'vip_port_id': '5bcbe2d7-994f-40de-87ab-07aa632f0133', 'provider': None, 'pools': [], 'id': 'e2d277f7-eca2-46a4-bf2c-655856fd8733', 'operating_status': 'OFFLINE', 'name': 'JL-B01-POD1-CORE-LB-7'}, {'subnets': {u'd79ef712-c1e3-4860-9343-d1702b9976aa': {'description': u'', 'tags': [], 'updated_at': '2020-09-25T05:29:56Z', 'ipv6_ra_mode': None, 'allocation_pools': [{'start': u'10.230.44.2', 'end': u'10.230.44.30'}], 'host_routes': [], 'revision_number': 1, 'ipv6_address_mode': None, 'cidr': u'10.230.44.0/27', 'id': u'd79ef712-c1e3-4860-9343-d1702b9976aa', 'subnetpool_id': None, 'service_types': [], 'available_ips': [{'start': '10.230.44.3', 'end': '10.230.44.3'}, {'start': '10.230.44.10', 'end': '10.230.44.12'}, {'start': '10.230.44.14', 'end': '10.230.44.14'}, {'start': '10.230.44.16', 'end': '10.230.44.17'}, {'start': '10.230.44.19', 'end': '10.230.44.19'}, {'start': '10.230.44.21', 'end': '10.230.44.21'}, {'start': '10.230.44.23', 'end': '10.230.44.25'}, {'start': '10.230.44.28', 'end': '10.230.44.30'}], 'name': u'LB-VIP', 'enable_dhcp': True, 'network_id': u'7801d370-530c-4c81-8a9f-c0b499dda220', 'tenant_id': u'62c38230485b4794a8eedece5dac9192', 'created_at': '2020-09-25T02:25:44Z', 'dns_nameservers': [], 'available_ip_number': 15, 'gateway_ip': u'10.230.44.1', 'ip_version': 4, 'shared': False, 'project_id': u'62c38230485b4794a8eedece5dac9192'}}, 'listeners': [], 'healthmonitors': [], 'members': [], 'l7policy_rules': [], 'pools': [], 'l7policies': [], 'networks': {u'7801d370-530c-4c81-8a9f-c0b499dda220': {'provider:physical_network': u'f5network1', 'updated_at': '2020-09-25T05:29:56Z', 'revision_number': 5, 'provider:network_type': u'vlan', 'id': u'7801d370-530c-4c81-8a9f-c0b499dda220', 'router:external': False, 'availability_zone_hints': [], 'availability_zones': [], 'provider:segmentation_id': 3020, 'ipv4_address_scope': None, 'shared': False, 'project_id': u'62c38230485b4794a8eedece5dac9192', 'status': u'ACTIVE', 'subnets': [u'd79ef712-c1e3-4860-9343-d1702b9976aa'], 'description': u'', 'tags': [], 'ipv6_address_scope': None, 'qos_policy_id': None, 'name': u'LB-VIP', 'admin_state_up': True, 'tenant_id': u'62c38230485b4794a8eedece5dac9192', 'created_at': '2020-09-25T02:18:35Z', 'mtu': 1450, 'vlan_transparent': False}}, 'loadbalancer': {'availability_zone_hints': [], 'description': '', 'admin_state_up': True, 'network_id': u'7801d370-530c-4c81-8a9f-c0b499dda220', 'tenant_id': '62c38230485b4794a8eedece5dac9192', 'provisioning_status': 'PENDING_CREATE', 'listeners': [], 'vip_port': {'status': u'DOWN', 'binding:host_id': u'POD1_CORE3', 'description': None, 'allowed_address_pairs': [], 'tags': [], 'extra_dhcp_opts': [], 'updated_at': '2020-10-05T02:20:26Z', 'device_owner': u'network:f5lbaasv2', 'revision_number': 7, 'binding:profile': {}, 'fixed_ips': [{'subnet_id': u'd79ef712-c1e3-4860-9343-d1702b9976aa', 'ip_address': u'10.230.44.15'}], 'id': u'5bcbe2d7-994f-40de-87ab-07aa632f0133', 'security_groups': [u'529d0bd0-0b3d-4e4f-941e-74cba6273c8e'], 'device_id': u'e2d277f7-eca2-46a4-bf2c-655856fd8733', 'name': u'loadbalancer-e2d277f7-eca2-46a4-bf2c-655856fd8733', 'admin_state_up': True, 'network_id': u'7801d370-530c-4c81-8a9f-c0b499dda220', 'tenant_id': u'62c38230485b4794a8eedece5dac9192', 'binding:vif_details': {}, 'binding:vnic_type': u'baremetal', 'binding:vif_type': u'other', 'qos_policy_id': None, 'mac_address': u'fa:16:3e:cb:0b:27', 'project_id': u'62c38230485b4794a8eedece5dac9192', 'created_at': '2020-10-05T02:20:16Z'}, 'vip_subnet_id': 'd79ef712-c1e3-4860-9343-d1702b9976aa', 'vip_address': '10.230.44.15', 'vip_port_id': '5bcbe2d7-994f-40de-87ab-07aa632f0133', 'provider': None, 'pools': [], 'id': 'e2d277f7-eca2-46a4-bf2c-655856fd8733', 'operating_status': 'OFFLINE', 'name': 'JL-B01-POD1-CORE-LB-7'}}, u'POD1_CORE3') {} wrapper /usr/lib/python2.7/site-packages/oslo_log/helpers.py:66`,
		"call_f5driver": 
			`2020-10-05 10:20:17.251 117825 DEBUG f5lbaasdriver.v2.bigip.driver_v2 [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - default default] f5lbaasdriver.v2.bigip.driver_v2.LoadBalancerManager method create called with arguments (<neutron_lib.context.Context object at 0x284cb250>, <neutron_lbaas.services.loadbalancer.data_models.LoadBalancer object at 0xdb44250>) {} wrapper /usr/lib/python2.7/site-packages/oslo_log/helpers.py:66`,
		"neutron_api_v2_base": 
			`2020-10-05 10:20:15.791 117825 DEBUG neutron.api.v2.base [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - default default] Request body: {u'loadbalancer': {u'vip_subnet_id': u'd79ef712-c1e3-4860-9343-d1702b9976aa', u'provider': u'core', u'name': u'JL-B01-POD1-CORE-LB-7', u'admin_state_up': True}} prepare_request_body /usr/lib/python2.7/site-packages/neutron/api/v2/base.py:713`,
		"call_f5agent":
			`2020-10-05 10:19:16.315 295263 DEBUG f5_openstack_agent.lbaasv2.drivers.bigip.agent_manager [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - - -] f5_openstack_agent.lbaasv2.drivers.bigip.agent_manager.LbaasAgentManager method create_loadbalancer called with arguments (<neutron_lib.context.Context object at 0x7351290>,) {u'service': {u'subnets': {u'd79ef712-c1e3-4860-9343-d1702b9976aa': {u'updated_at': u'2020-09-25T05:29:56Z', u'ipv6_ra_mode': None, u'allocation_pools': [{u'start': u'10.230.44.2', u'end': u'10.230.44.30'}], u'host_routes': [], u'revision_number': 1, u'ipv6_address_mode': None, u'id': u'd79ef712-c1e3-4860-9343-d1702b9976aa', u'available_ips': [{u'start': u'10.230.44.3', u'end': u'10.230.44.3'}, {u'start': u'10.230.44.10', u'end': u'10.230.44.12'}, {u'start': u'10.230.44.14', u'end': u'10.230.44.14'}, {u'start': u'10.230.44.16', u'end': u'10.230.44.17'}, {u'start': u'10.230.44.19', u'end': u'10.230.44.19'}, {u'start': u'10.230.44.21', u'end': u'10.230.44.21'}, {u'start': u'10.230.44.23', u'end': u'10.230.44.25'}, {u'start': u'10.230.44.28', u'end': u'10.230.44.30'}], u'dns_nameservers': [], u'gateway_ip': u'10.230.44.1', u'shared': False, u'project_id': u'62c38230485b4794a8eedece5dac9192', u'description': u'', u'tags': [], u'available_ip_number': 15, u'cidr': u'10.230.44.0/27', u'subnetpool_id': None, u'service_types': [], u'name': u'LB-VIP', u'enable_dhcp': True, u'network_id': u'7801d370-530c-4c81-8a9f-c0b499dda220', u'tenant_id': u'62c38230485b4794a8eedece5dac9192', u'created_at': u'2020-09-25T02:25:44Z', u'ip_version': 4}}, u'listeners': [], u'healthmonitors': [], u'members': [], u'l7policy_rules': [], u'pools': [], u'l7policies': [], u'networks': {u'7801d370-530c-4c81-8a9f-c0b499dda220': {u'provider:physical_network': u'f5network1', u'updated_at': u'2020-09-25T05:29:56Z', u'revision_number': 5, u'mtu': 1450, u'id': u'7801d370-530c-4c81-8a9f-c0b499dda220', u'router:external': False, u'availability_zone_hints': [], u'availability_zones': [], u'provider:segmentation_id': 3020, u'ipv4_address_scope': None, u'shared': False, u'project_id': u'62c38230485b4794a8eedece5dac9192', u'status': u'ACTIVE', u'subnets': [u'd79ef712-c1e3-4860-9343-d1702b9976aa'], u'description': u'', u'tags': [], u'ipv6_address_scope': None, u'qos_policy_id': None, u'name': u'LB-VIP', u'admin_state_up': True, u'tenant_id': u'62c38230485b4794a8eedece5dac9192', u'created_at': u'2020-09-25T02:18:35Z', u'provider:network_type': u'vlan', u'vlan_transparent': False}}, u'loadbalancer': {u'availability_zone_hints': [], u'description': u'', u'provisioning_status': u'PENDING_CREATE', u'network_id': u'7801d370-530c-4c81-8a9f-c0b499dda220', u'tenant_id': u'62c38230485b4794a8eedece5dac9192', u'admin_state_up': True, u'provider': None, u'id': u'e2d277f7-eca2-46a4-bf2c-655856fd8733', u'pools': [], u'listeners': [], u'vip_port_id': u'5bcbe2d7-994f-40de-87ab-07aa632f0133', u'vip_address': u'10.230.44.15', u'vip_subnet_id': u'd79ef712-c1e3-4860-9343-d1702b9976aa', u'vip_port': {u'allowed_address_pairs': [], u'extra_dhcp_opts': [], u'updated_at': u'2020-10-05T02:20:26Z', u'device_owner': u'network:f5lbaasv2', u'revision_number': 7, u'binding:profile': {}, u'fixed_ips': [{u'subnet_id': u'd79ef712-c1e3-4860-9343-d1702b9976aa', u'ip_address': u'10.230.44.15'}], u'id': u'5bcbe2d7-994f-40de-87ab-07aa632f0133', u'security_groups': [u'529d0bd0-0b3d-4e4f-941e-74cba6273c8e'], u'binding:vif_details': {}, u'binding:vif_type': u'other', u'mac_address': u'fa:16:3e:cb:0b:27', u'device_id': u'e2d277f7-eca2-46a4-bf2c-655856fd8733', u'status': u'DOWN', u'binding:host_id': u'POD1_CORE3', u'description': None, u'tags': [], u'qos_policy_id': None, u'project_id': u'62c38230485b4794a8eedece5dac9192', u'name': u'loadbalancer-e2d277f7-eca2-46a4-bf2c-655856fd8733', u'admin_state_up': True, u'network_id': u'7801d370-530c-4c81-8a9f-c0b499dda220', u'tenant_id': u'62c38230485b4794a8eedece5dac9192', u'created_at': u'2020-10-05T02:20:16Z', u'binding:vnic_type': u'baremetal'}, u'operating_status': u'OFFLINE', u'name': u'JL-B01-POD1-CORE-LB-7'}}, u'loadbalancer': {u'availability_zone_hints': [], u'description': u'', u'provisioning_status': u'PENDING_CREATE', u'tenant_id': u'62c38230485b4794a8eedece5dac9192', u'admin_state_up': True, u'provider': None, u'pools': [], u'listeners': [], u'vip_port_id': u'5bcbe2d7-994f-40de-87ab-07aa632f0133', u'vip_address': u'10.230.44.15', u'vip_subnet_id': u'd79ef712-c1e3-4860-9343-d1702b9976aa', u'id': u'e2d277f7-eca2-46a4-bf2c-655856fd8733', u'operating_status': u'OFFLINE', u'name': u'JL-B01-POD1-CORE-LB-7'}} wrapper /usr/lib/python2.7/site-packages/oslo_log/helpers.py:66`,
		"rest_call_bigip":
			`2020-10-05 10:19:16.317 295263 DEBUG root [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - - -] get WITH uri: https://10.216.177.8:443/mgmt/tm/sys/folder/~CORE_62c38230485b4794a8eedece5dac9192 AND suffix:  AND kwargs: {} wrapper /usr/lib/python2.7/site-packages/icontrol/session.py:257`,
		"report_loadbalancer_status":
			`2020-10-05 10:19:18.411 295263 DEBUG f5_openstack_agent.lbaasv2.drivers.bigip.plugin_rpc [req-92db71fb-8513-431b-ac79-5423a749b6d7 009ac6496334436a8eba8daa510ef659 62c38230485b4794a8eedece5dac9192 - - -] f5_openstack_agent.lbaasv2.drivers.bigip.plugin_rpc.LBaaSv2PluginRPC method update_loadbalancer_status called with arguments (u'e2d277f7-eca2-46a4-bf2c-655856fd8733', 'ACTIVE', 'ONLINE', u'JL-B01-POD1-CORE-LB-7') {} wrapper /usr/lib/python2.7/site-packages/oslo_log/helpers.py:66`,
	}
}
