/*
Baserow API spec

For more information about our REST API, please visit [this page](https://baserow.io/docs/apis%2Frest-api).  For more information about our deprecation policy, please visit [this page](https://baserow.io/docs/apis%2Fdeprecations).

API version: 1.18.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"fmt"
)

// FilterActionTypeEnum * `create_group` - create_group * `delete_group` - delete_group * `update_group` - update_group * `order_groups` - order_groups * `create_application` - create_application * `update_application` - update_application * `delete_application` - delete_application * `order_applications` - order_applications * `duplicate_application` - duplicate_application * `install_template` - install_template * `create_group_invitation` - create_group_invitation * `delete_group_invitation` - delete_group_invitation * `accept_group_invitation` - accept_group_invitation * `reject_group_invitation` - reject_group_invitation * `update_group_invitation_permissions` - update_group_invitation_permissions * `leave_group` - leave_group * `create_snapshot` - create_snapshot * `delete_snapshot` - delete_snapshot * `restore_snapshot` - restore_snapshot * `empty_trash` - empty_trash * `restore_from_trash` - restore_from_trash * `create_user` - create_user * `update_user` - update_user * `schedule_user_deletion` - schedule_user_deletion * `cancel_user_deletion` - cancel_user_deletion * `sign_in_user` - sign_in_user * `change_user_password` - change_user_password * `send_reset_user_password` - send_reset_user_password * `reset_user_password` - reset_user_password * `create_db_token` - create_db_token * `update_db_token_name` - update_db_token_name * `update_db_token_permissions` - update_db_token_permissions * `rotate_db_token_key` - rotate_db_token_key * `delete_db_token_key` - delete_db_token_key * `create_webhook` - create_webhook * `delete_webhook` - delete_webhook * `update_webhook` - update_webhook * `export_table` - export_table * `import_database_from_airtable` - import_database_from_airtable * `create_table` - create_table * `delete_table` - delete_table * `order_tables` - order_tables * `update_table` - update_table * `duplicate_table` - duplicate_table * `create_row` - create_row * `create_rows` - create_rows * `import_rows` - import_rows * `delete_row` - delete_row * `delete_rows` - delete_rows * `move_row` - move_row * `update_row` - update_row * `update_rows` - update_rows * `create_view` - create_view * `duplicate_view` - duplicate_view * `delete_view` - delete_view * `order_views` - order_views * `update_view` - update_view * `create_view_filter` - create_view_filter * `update_view_filter` - update_view_filter * `delete_view_filter` - delete_view_filter * `create_view_sort` - create_view_sort * `update_view_sort` - update_view_sort * `delete_view_sort` - delete_view_sort * `rotate_view_slug` - rotate_view_slug * `update_view_field_options` - update_view_field_options * `create_decoration` - create_decoration * `update_decoration` - update_decoration * `delete_decoration` - delete_decoration * `create_field` - create_field * `delete_field` - delete_field * `update_field` - update_field * `duplicate_field` - duplicate_field * `create_row_comment` - create_row_comment * `create_team` - create_team * `update_team` - update_team * `delete_team` - delete_team * `create_team_subject` - create_team_subject * `delete_team_subject` - delete_team_subject * `batch_assign_role` - batch_assign_role
type FilterActionTypeEnum string

// List of FilterActionTypeEnum

// All allowed values of FilterActionTypeEnum enum
var AllowedFilterActionTypeEnumEnumValues = []FilterActionTypeEnum{
	"create_group",
	"delete_group",
	"update_group",
	"order_groups",
	"create_application",
	"update_application",
	"delete_application",
	"order_applications",
	"duplicate_application",
	"install_template",
	"create_group_invitation",
	"delete_group_invitation",
	"accept_group_invitation",
	"reject_group_invitation",
	"update_group_invitation_permissions",
	"leave_group",
	"create_snapshot",
	"delete_snapshot",
	"restore_snapshot",
	"empty_trash",
	"restore_from_trash",
	"create_user",
	"update_user",
	"schedule_user_deletion",
	"cancel_user_deletion",
	"sign_in_user",
	"change_user_password",
	"send_reset_user_password",
	"reset_user_password",
	"create_db_token",
	"update_db_token_name",
	"update_db_token_permissions",
	"rotate_db_token_key",
	"delete_db_token_key",
	"create_webhook",
	"delete_webhook",
	"update_webhook",
	"export_table",
	"import_database_from_airtable",
	"create_table",
	"delete_table",
	"order_tables",
	"update_table",
	"duplicate_table",
	"create_row",
	"create_rows",
	"import_rows",
	"delete_row",
	"delete_rows",
	"move_row",
	"update_row",
	"update_rows",
	"create_view",
	"duplicate_view",
	"delete_view",
	"order_views",
	"update_view",
	"create_view_filter",
	"update_view_filter",
	"delete_view_filter",
	"create_view_sort",
	"update_view_sort",
	"delete_view_sort",
	"rotate_view_slug",
	"update_view_field_options",
	"create_decoration",
	"update_decoration",
	"delete_decoration",
	"create_field",
	"delete_field",
	"update_field",
	"duplicate_field",
	"create_row_comment",
	"create_team",
	"update_team",
	"delete_team",
	"create_team_subject",
	"delete_team_subject",
	"batch_assign_role",
}

func (v *FilterActionTypeEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := FilterActionTypeEnum(value)
	for _, existing := range AllowedFilterActionTypeEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid FilterActionTypeEnum", value)
}

// NewFilterActionTypeEnumFromValue returns a pointer to a valid FilterActionTypeEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewFilterActionTypeEnumFromValue(v string) (*FilterActionTypeEnum, error) {
	ev := FilterActionTypeEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for FilterActionTypeEnum: valid values are %v", v, AllowedFilterActionTypeEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v FilterActionTypeEnum) IsValid() bool {
	for _, existing := range AllowedFilterActionTypeEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to FilterActionTypeEnum value
func (v FilterActionTypeEnum) Ptr() *FilterActionTypeEnum {
	return &v
}

type NullableFilterActionTypeEnum struct {
	value *FilterActionTypeEnum
	isSet bool
}

func (v NullableFilterActionTypeEnum) Get() *FilterActionTypeEnum {
	return v.value
}

func (v *NullableFilterActionTypeEnum) Set(val *FilterActionTypeEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableFilterActionTypeEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableFilterActionTypeEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFilterActionTypeEnum(val *FilterActionTypeEnum) *NullableFilterActionTypeEnum {
	return &NullableFilterActionTypeEnum{value: val, isSet: true}
}

func (v NullableFilterActionTypeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFilterActionTypeEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

