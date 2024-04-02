scheme funeral_assistance
    label Funeral Assistance
    criteria occupation
        table fm
        column occupation
        operator equals
        value Construction Worker
    criteria health_status
        table fm
        column health_status
        operator equals
        value Deceased
    criteria card_registration
        table fm
        column bocw_card_registration_date
        operator age_gte // generalising age to be difference between the column and today
        granularity year
        value 0

    evaluation occupation && health_status && card_registration