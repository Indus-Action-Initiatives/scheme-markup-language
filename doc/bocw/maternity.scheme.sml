scheme maternity_benefit
    label Maternity Benefit
    description """Maternity benefits of Rs. 30,000 to registered women members 
and wives of male members (upto 2 children). (Rule – 271) – from the date of joining membership of the fund."""
    criteria gender_and_age
        combine OR
            term male_gender_and_age
                combine AND
                    term male_gender
                        table fm
                        column gender
                        operator equals
                        value male
                    term male_age
                        table fm
                        column age
                        operator gte
                        value 21   
            term female_gender_and_age         
                combine AND
                    term female_gender
                        table fm
                        column gender
                        operator equals
                        value female
                    term female_age
                        table fm
                        column age
                        operator gte
                        value 18
    criteria occupation
        table fm
        column occupation
        operator equals
        value Construction Worker
    criteria marital_status
        table fm
        column marital_status
        operator equals
        value Married
    criteria pregnancy_status
        table fm
        column pregnancy_status
        operator IN
        value ['Delivered first child', 'Delivered second child']
    criteria number_of_children
        table fm
        column number_of_children
        operator gte
        value 1
    criteria card_registration
        table fm
        column bocw_card_registration_date
        operator age_gte // generalising age to be difference between the column and today
        granularity year
        value 0

    evaluation age && occupation && marital_status && pregnancy_status && number_of_children && card_registration