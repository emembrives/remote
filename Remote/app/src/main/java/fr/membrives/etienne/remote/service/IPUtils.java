package fr.membrives.etienne.remote.service;

import android.util.Log;

import java.math.BigInteger;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.HashSet;
import java.util.Set;

/**
 * Created by etienne on 23/05/15.
 */
public class IPUtils {
    public static Set<String> getIpsInSubnet(int ip, int netmask, int index) {
        if (index == 32) {
            Set<String> set = new HashSet<>();
            byte localComponent = BigInteger.valueOf(ip).toByteArray()[3];
            if (localComponent != (byte) 255) {
                // Remove the broadcast address
                set.add(ipToString(ip));
            }
            return set;
        }

        if (((netmask >> 31 - index) & 1) == 1) {
            return getIpsInSubnet(ip, netmask, index + 1);
        }

        int zeroIp = ip & ~(1 << (31 - index));
        Set<String> zeroIps = getIpsInSubnet(zeroIp, netmask, index + 1);
        int oneIp = ip | (1 << (31 - index));
        Set<String> oneIps = getIpsInSubnet(oneIp, netmask, index + 1);
        zeroIps.addAll(oneIps);
        return zeroIps;
    }

    private static String ipToString(int ipAddress) {
        byte[] ipByteArray = BigInteger.valueOf(ipAddress).toByteArray();

        String ipAddressString;
        try {
            ipAddressString = InetAddress.getByAddress(ipByteArray).getHostAddress();
        } catch (UnknownHostException e) {
            throw new RuntimeException("IP address: " + Integer.toString(ipAddress), e);
        }
        return ipAddressString;
    }
}
