scheme work_tools_grant
    label Grant for purchase of work-related tools
    criteria age
        table fm
        column age
        operator lt
        value 55
    criteria occupation
        table fm
        column occupation
        operator equals
        value Construction Worker
    criteria card_registration
        table fm
        column bocw_card_registration_date
        operator age_gte // generalising age to be difference between the column and today
        granularity year
        value 3

    evaluation age && occupation && card_registration