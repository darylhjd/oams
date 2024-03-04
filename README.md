<div align="center">

# OAMS 📑
Online Attendance Management System for NTU.

OAMS comes with a [website](https://oams-webserver-staging.jollyflower-f772283d.southeastasia.azurecontainerapps.io/)! If you are an administrative user, you may login to explore more features.

**<i>NOTE: OAMS is currently in development and is not open for public testing.
If you are an NTU student and wish to try out this system, please feel free to drop me a message.</i>**

</div>
<div align="justify">

## External Services
Although OAMS comes with a built-in signature authentication system for attendance taking (which is used by the website), it is not the only possible method for attendance authentication.

OAMS provides avenues for integration with external services. External services are programmatic services that act as middlemen, allowing them to introduce different methods of attendance authentication before updating
a user's attendance. For example, biometric, QR-Code, or geo-fencing capabilities for authenticating a user's identity can be added before sending a request to change a user's attendance record.

Anyone is capable of creating an external service if they are so inclined!

### API Endpoints
The following endpoints are available for usage by external services.

<table>
  <tr>
    <th>This is a header</th>
  </tr>
</table>
