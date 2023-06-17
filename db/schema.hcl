schema "public" {}

table "teams" {
    schema = schema.public

    column "uuid" {
        type = uuid
        default = sql("gen_random_uuid()")
    }
    column "name" {
        type = text
    }
    column "description" {
        type = text
        null = true
    }
    column "created_at" {
        type = timestamptz
    }
    column "updated_at" {
        type = timestamptz
    }
    column "deleted_at" {
        type = timestamptz
        null = true
    }

    primary_key {
        columns = [column.uuid]
    }
}

table "team_members" {
    schema = schema.public

    column "uuid" {
        type = uuid
        default = sql("gen_random_uuid()")
    }
    column "team_uuid" {
        type = uuid
    }
    column "email" {
        type = text
    }
    column "is_admin" {
        type = boolean
        default = false
    }
    column "is_owner" {
        type = boolean
        default = false
    }
    column "created_at" {
        type = timestamptz
    }
    column "updated_at" {
        type = timestamptz
    }
    column "deleted_at" {
        type = timestamptz
        null = true
    }

    primary_key {
        columns = [column.uuid]
    }
    foreign_key "team_uuid" {
        columns = [column.team_uuid]
        ref_columns = [table.teams.column.uuid]
        on_update = NO_ACTION
        on_delete = NO_ACTION
    }
}