# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/admin/admin.proto](#api/admin/admin.proto)
  
  
  
    - [ZtpAdminService](#admin.ZtpAdminService)
  

- [api/admin/roles.proto](#api/admin/roles.proto)
    - [DeviceConfig](#admin.DeviceConfig)
    - [DevicePipeline](#admin.DevicePipeline)
    - [DeviceProperty](#admin.DeviceProperty)
    - [DeviceRoleChange](#admin.DeviceRoleChange)
    - [DeviceRoleChangeRequest](#admin.DeviceRoleChangeRequest)
    - [DeviceRoleChangeResponse](#admin.DeviceRoleChangeResponse)
    - [DeviceRoleConfig](#admin.DeviceRoleConfig)
    - [DeviceRoleRequest](#admin.DeviceRoleRequest)
  
    - [DeviceRoleChange.ChangeType](#admin.DeviceRoleChange.ChangeType)
    - [DeviceRoleChangeRequest.ChangeType](#admin.DeviceRoleChangeRequest.ChangeType)
  
  
    - [DeviceRoleService](#admin.DeviceRoleService)
  

- [Scalar Value Types](#scalar-value-types)



<a name="api/admin/admin.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/admin/admin.proto


 

 

 


<a name="admin.ZtpAdminService"></a>

### ZtpAdminService
ZtpAdminService provides means for enhanced interactions with the zero-touch-provisioning subsystem.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|

 



<a name="api/admin/roles.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/admin/roles.proto



<a name="admin.DeviceConfig"></a>

### DeviceConfig
DeviceConfig is a set of initial configuration properties to be applied to a device.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| softwareVersion | [string](#string) |  |  |
| properties | [DeviceProperty](#admin.DeviceProperty) | repeated |  |






<a name="admin.DevicePipeline"></a>

### DevicePipeline
DevicePipeline carries information about the P4 pipeline configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pipeconf | [string](#string) |  |  |
| driver | [string](#string) |  |  |






<a name="admin.DeviceProperty"></a>

### DeviceProperty
DeviceProperty is a path/type/value tuple


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |
| type | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="admin.DeviceRoleChange"></a>

### DeviceRoleChange
DeviceRoleChange is an event describing a change to a device role configuration.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| change | [DeviceRoleChange.ChangeType](#admin.DeviceRoleChange.ChangeType) |  |  |
| config | [DeviceRoleConfig](#admin.DeviceRoleConfig) |  |  |






<a name="admin.DeviceRoleChangeRequest"></a>

### DeviceRoleChangeRequest
DeviceRoleChangeRequest is a request for a change to a device role configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| change | [DeviceRoleChangeRequest.ChangeType](#admin.DeviceRoleChangeRequest.ChangeType) |  |  |
| config | [DeviceRoleConfig](#admin.DeviceRoleConfig) |  |  |






<a name="admin.DeviceRoleChangeResponse"></a>

### DeviceRoleChangeResponse
DeviceRoleChangeResponse is a response for a change to a device role configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| change | [DeviceRoleChange](#admin.DeviceRoleChange) |  |  |






<a name="admin.DeviceRoleConfig"></a>

### DeviceRoleConfig
DeviceRoleConfig carries the template configuration associated with a device role


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role | [string](#string) |  |  |
| config | [DeviceConfig](#admin.DeviceConfig) |  |  |
| pipeline | [DevicePipeline](#admin.DevicePipeline) |  |  |






<a name="admin.DeviceRoleRequest"></a>

### DeviceRoleRequest
DeviceRoleRequest is a request for device role configuration.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role | [string](#string) |  |  |





 


<a name="admin.DeviceRoleChange.ChangeType"></a>

### DeviceRoleChange.ChangeType


| Name | Number | Description |
| ---- | ------ | ----------- |
| UPDATED | 0 |  |
| ADDED | 1 |  |
| DELETED | 2 |  |



<a name="admin.DeviceRoleChangeRequest.ChangeType"></a>

### DeviceRoleChangeRequest.ChangeType


| Name | Number | Description |
| ---- | ------ | ----------- |
| UPDATE | 0 |  |
| ADD | 1 |  |
| DELETE | 2 |  |


 

 


<a name="admin.DeviceRoleService"></a>

### DeviceRoleService
DeviceRoleService provides means for setting up device role configurations
in support of zero-touch provisioning activities.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Set | [DeviceRoleChangeRequest](#admin.DeviceRoleChangeRequest) | [DeviceRoleChangeResponse](#admin.DeviceRoleChangeResponse) | Set provides means to add, update or delete device role configuration. |
| Get | [DeviceRoleRequest](#admin.DeviceRoleRequest) | [DeviceRoleConfig](#admin.DeviceRoleConfig) stream | Get provides means to query device role configuration. |
| Subscribe | [DeviceRoleRequest](#admin.DeviceRoleRequest) | [DeviceRoleChange](#admin.DeviceRoleChange) stream | Subscribe provides means to monitor changes in the device role configuration. |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

