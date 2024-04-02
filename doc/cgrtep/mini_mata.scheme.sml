scheme mini_mata
    label Mini Mata Mahtari Jatan Yojna
    description """Female construction worker gets INR 20,000 during pregnancy"""
    criteria gender
        column gender
        table fm
        operator IN
        value ['female', 'other']
    criteria pregnancy
        column pregnancy_status
        table fm
        operator equals
        value true
    criteria bocw_card
        column has_bocw_card
        table fm
        operator equals
        value true
    criteria bocw_card_issue_date
        column bocw_card_issue_date
        table fm
        operator age_gte
        value 90
        granularity day
    
    evaluation gender && pregnancy && bocw_card && bocw_card_issue_date
        