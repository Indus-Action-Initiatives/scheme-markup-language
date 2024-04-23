scheme shramik
    label Chief Minister Shramik Tool Assistance Scheme
    description """A one-time benefit of Rs. 1000 is provided to buy a tool kit (if tool kit is available then the kit is given)"""    
    criteria bocw_card
        column has_bocw_card
        table fm
        operator equals
        value true
    criteria age
        column dob
        table fm
        operator age_between
        value [18, 50]
        granularity year
    
    evaluation bocw_card && age
        