# forwarding_rule.py
# : Define 3 different forwarding rules
# :     - select ip from server_list
# :     (1) round-robin
# :     (2) random_choice
# :     (3) select_least_count_ip
# --------------------------------------------------------------------
import random



# Select ip(round_robin) from server list with specific (protocol, port)
def round_robin(protocol, port):
    key = (protocol.lower(), port)
    from LB.table.server_list import server_list

    if key in server_list:

        data_dict_entry = server_list[key]

        ips_list = data_dict_entry.get("ips", [])

        # get index for round-robin (stored in server list)
        round_robin_index = data_dict_entry.get("round_robin")
        if round_robin_index is None:
            return None

        if ips_list:
            index = (round_robin_index + 1) % len(ips_list)
            selected_ip = ips_list[index].get("ip", None)

            # update index for round-robin
            data_dict_entry["round_robin"] = index
            return selected_ip
    return None


# Select ip(random_choice) from server list with specific (protocol, port)
def random_choice(protocol, port):
    key = (protocol.lower(), port)
    from LB.table.server_list import server_list
    if key in server_list:
        data_dict_entry = server_list[key]

        ips_list = data_dict_entry.get("ips", [])

        if ips_list:
            selected_ip = random.choice(ips_list).get("ip", None)
            return selected_ip
    return None


# Select ip(with least user_count) from server list with specific (protocol, port)
def select_least_count_ip(protocol, port):
    key = (protocol.lower(), port)
    from LB.table.server_list import server_list

    if key in server_list:
        data_dict = server_list[key]

        # "ips" 키를 가진 리스트 추출
        ips_list = data_dict.get("ips", [])

        # get ip with minimum user_count (stored in server list)
        if ips_list:
            min_count_data = min(ips_list, key=lambda x: x.get("user_count", float('inf')))
            return min_count_data.get("ip", None)
    return None
