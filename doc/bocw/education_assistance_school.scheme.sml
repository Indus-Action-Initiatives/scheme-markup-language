scheme education_assistance_school
    label Education Assistance (School)
    criteria occupation
        table fm
        column occupation
        operator equals
        value Construction Worker
    criteria number_of_children
        table fm
        column number_of_children
        operator gte
        value 1
    criteria children_in_school
        table fm
        column children_school_or_college
        operator equals
        value school
    criteria card_registration
        table fm
        column bocw_card_registration_date
        operator age_gte // generalising age to be difference between the column and today
        granularity year
        value 0

    evaluation occupation && number_of_children && children_in_school && card_registration