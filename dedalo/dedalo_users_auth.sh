#!/bin/bash

users_auth () {
	. /opt/icaro/dedalo/config
	HS_DIGEST=$(echo -n "$HS_SECRET$HS_UUID" | md5sum | awk '{print $1}')

	devices_list=$(dedalo query list | grep -e dnat -e temporary )

	if [ $? -ne 0 ]; then
		#No device found
		return 0
	fi

	IFS=$'\n'
	for device in $devices_list; do

		device_sessionid=$(echo $device | awk '{print $4}')

		resp_body=$(curl -L -f -s "${HS_AAA_URL}/auth?digest=${HS_DIGEST}&uuid=${HS_UUID}&secret=${HS_SECRET}&sessionid=${device_sessionid}")

		if [ $? -eq 0 ]; then

			auth_type=$(echo ${resp_body} | jq -r '.type')
			device_mac=$(echo $device | awk '{print $1}')
			device_username=$(echo $device | awk '{print $6}')

			if [ "$auth_type" = "login" ]; then

				username=$(echo ${resp_body} | jq -r '.username')
				password=$(echo ${resp_body} | jq -r '.password')

				#In case of temporary session, first deauthenticate the device
				if [ "$device_username" = "temporary" ]; then
					dedalo query logout mac ${device_mac}
				fi

				#Authenticate the device
				dedalo query login mac ${device_mac} username "${username}" password "${password}"

			elif [ "$auth_type" = "temporary" ]; then

				#Check if a temporary session does not already exist
				if [ "$device_username" != "temporary" ]; then

					sessiontimeout=$(echo ${resp_body} | jq -r '.sessiontimeout')
					#Open a temporary session
					dedalo query authorize mac ${device_mac} username "temporary" sessiontimeout ${sessiontimeout}

				fi
			fi
		fi
done
}

while true
do
	users_auth
	sleep 10s
done
