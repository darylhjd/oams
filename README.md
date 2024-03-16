<div align="center">

# OAMS 📑

Online Attendance Management System for NTU.

OAMS comes with a [website](https://oams-webserver-staging.jollyflower-f772283d.southeastasia.azurecontainerapps.io/)!
If you are an administrative user, you may login to explore more features.

**<i>NOTE: OAMS is currently in development and is not open for public testing.
If you are an NTU student and wish to try out this system, please feel free to drop me a message.</i>**

</div>
<div align="justify">

## External Services

Although OAMS comes with a built-in signature authentication system for attendance taking (which is used by the
website), it is not the only possible method for attendance authentication.

OAMS provides avenues for integration with external services. External services act as middlemen, allowing them to
introduce different methods of attendance authentication before updating a user's attendance. For example, biometric,
QR-Code, or geo-fencing capabilities for authenticating a user's identity can be added before sending a request to
change a user's attendance record.

Anyone is capable of creating an external service if they are so inclined!

### API Endpoints

The following endpoints are available for usage by external services.

</div>

<table>
    <tr>
        <th>Endpoint</th>
        <th>Description</th>
        <th>Methods</th>
        <th>Codes</th>
        <th>Request</th>
        <th>Result</th>
    </tr>
    <tr />
    <tr>
        <td>/upcoming-class-group-sessions</td>
        <td>Gets upcoming class group sessions.</td>
        <td>GET</td>
        <td>200: Success</td>
        <td>-</td>
        <td>
            <details>
            <summary>Response</summary>
            <pre>
<code>{
    "result": boolean,
    "upcoming_class_group_session": []UpcomingClassGroupSession
}</code>
            </pre>
            </details>
        </td>
    </tr>
    <tr />
    <tr>
        <td rowspan="2">/upcoming-class-group-sessions/{sessionId}/attendances</td>
        <td rowspan="2">Get attendance records for class group session with {sessionId}</td>
        <td rowspan="2">GET</td>
        <td rowspan="1">200: Success</td>
        <td rowspan="2">-</td>
        <td rowspan="2">
            <details>
            <summary>Response</summary>
            <pre>
<code>{
    "result": boolean,
    "upcoming_class_group_session": UpcomingClassGroupSession,
    "attendance_entries": []AttendanceEntry
}</code>
            </pre>
            </details>
        </td>
    </tr>
    <tr>
        <td>404: Session does not exist or is not upcoming.</td>
    </tr>
    <tr>
        <td rowspan="2">/upcoming-class-group-sessions/{sessionId}/attendances/{enrollmentId}</td>
        <td rowspan="2">Change attendance status for a given {sessionId} and {enrollmentId}.</td>
        <td rowspan="2">PATCH</td>
        <td rowspan="1">200: Success</td>
        <td rowspan="2">
            <details>
            <summary>Request</summary>
            <pre>
<code>{
    "attended": boolean
}</code>
            </pre>
            </details>
        </td>
        <td rowspan="2">
            <details>
            <summary>Response</summary>
            <pre>
<code>{
    "result": boolean,
    "attended: boolean
}</code>
            </pre>
            </details>
        </td>
    </tr>
    <tr>
        <td>401: Not allowed to change attendance</td>
    </tr>
</table>

<table>
    <tr>
        <th>Entity</th>
        <th>Format</th>
    </tr>
    <tr />
    <tr>
        <th>UpcomingClassGroupSession</th>
        <td>
            <details>
            <summary>Example</summary>
            <pre>
<code>{
    "id": 1,
    "start_time": "2024-01-15T08:30:00+08:00",
    "end_time": "2024-01-15T09:20:00+08:00",
    "venue": "TR+15 NORTH,NS4-05-93",
    "code": "SC1015",
    "year": 2023,
    "semester": "2",
    "name": "A21",
    "class_type": "TUT",
    "managing_role": "TEACHING_ASSISTANT" // Can be ignored for external services.
}</code>
            </pre>
            </details>
        </td>
    </tr>
    <tr />
    <tr>
        <th>AttendanceEntry</th>
        <td>
            <details>
            <summary>Example</summary>
            <pre>
<code>{
    "id": 2,
    "session_id": 1,
    "user_id": "TEST1345",
    "user_name": "JOHN TAN",
    "attended": false
}</code>
            </pre>
            </details>
        </td>
    </tr>
</table>
