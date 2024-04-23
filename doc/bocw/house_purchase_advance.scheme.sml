scheme house_purchase_advance
    label Advance for purchase or construction of the house
    criteria age
        table fm
        column age
        operator lt
        value 50
    criteria occupation
        table fm
        column occupation
        operator equals
        value Construction Worker
    criteria home_ownership_status
        table fm
        column home_ownership_status
        operator IN
        value ['Owner', 'Going to purchase']
    criteria card_registration
        table fm
        column bocw_card_registration_date
        operator age_gte // generalising age to be difference between the column and today
        granularity year
        value 5

    evaluation age && occupation && home_ownership_status && card_registration