env "local" {
    // schema definition
    src = "file://db/schema.hcl"

    // postgres instance
    // TODO: Read db params from an env file
    url = "postgres://manager:password@localhost:5432/team_app?search_path=public&sslmode=disable"

    // Define the URL of the Dev Database for this environment
    // See: https://atlasgo.io/concepts/dev-database
    dev = "docker://postgres/15"

    migration {
        // URL where the migration directory resides
        dir = "file://db/migrations"

        // Format of the migration directory
        // atlas (default) | flyway | liquibase | goose | golang-migrate | dbmate
        format = atlas
    }
}

env "test" {
    // schema definition
    src = "file://db/schema.hcl"

    // postgres instance
    // TODO: Read db params from an env file
    url = "postgres://test:test@localhost:5432/test?search_path=public&sslmode=disable"

    // Define the URL of the Dev Database for this environment
    // See: https://atlasgo.io/concepts/dev-database
    dev = "docker://postgres/15"

    migration {
        // URL where the migration directory resides
        dir = "file://db/migrations"

        // Format of the migration directory
        // atlas (default) | flyway | liquibase | goose | golang-migrate | dbmate
        format = atlas
    }
}