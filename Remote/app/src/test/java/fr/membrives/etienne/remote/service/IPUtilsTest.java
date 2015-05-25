package fr.membrives.etienne.remote.service;

import org.junit.Assert;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.robolectric.RobolectricGradleTestRunner;
import org.robolectric.annotation.Config;

import java.math.BigInteger;
import java.util.Set;

import fr.membrives.etienne.remote.BuildConfig;

import static org.hamcrest.CoreMatchers.hasItem;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.not;

@RunWith(RobolectricGradleTestRunner.class)
@Config(constants = BuildConfig.class)
public class IPUtilsTest {
    @Test
    public void testGetIpsInSubnet() throws Exception {
        byte[] ipAddress = {(byte) 192, (byte) 168, 1, 0};
        byte[] netmask = {(byte) 255, (byte) 255, (byte) 255, 0};
        Set<String> ips = IPUtils.getIpsInSubnet(new BigInteger(1, ipAddress).intValue(),
                new BigInteger(1, netmask).intValue(), 0);
        Assert.assertThat(255, is(ips.size()));
        Assert.assertThat(ips, hasItem("192.168.1.0"));
        Assert.assertThat(ips, hasItem("192.168.1.42"));
        Assert.assertThat(ips, hasItem("192.168.1.199"));
        Assert.assertThat(ips, not(hasItem("192.168.1.255")));
        Assert.assertThat(ips, not(hasItem("192.168.0.0")));
    }
}