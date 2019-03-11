<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <style>
        body {
            background-color: {{ .BgColor }};
            text-align: center;
            font-family: Lato, 'Helvetica Neue', Arial, Helvetica,
sans-serif;
        }

        h1 {
            color: white;
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
        }

        h2 {
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
            margin: 20px;
        }

        h3 {
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
        }

        h4 {
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
        }

        h5 {
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
        }

        p {
            color: red;
        }

        img {
            width: 150px;
        }

        .star {
            width: 50px;
            margin-bottom: 25px;
        }

        .title {
            margin: 20px;
        }

        .container {
            background: white;
            margin-left: 10px !important;
            margin-right: 10px !important;
            border-radius: .28571429rem;
            padding: 10px;
        }

        .divider {
            padding: 5px;
        }
    </style>
    <title>Review</title>
</head>

<body>
    <div class="title">
        <div class="container">
            <h2>Hai ricevuto un nuovo feedback!</h2>
            <h3>Questo è il messaggio dell'utente</h3>
            <div>
                {{ .Message }}
            </div>
            <h5>Questa è una mail automatica. Non rispondere</h5>
        </div>
        <div class="divider"></div>
        <div class="container">
            <h2>You have received new feedback!</h2>
            <h3>This is the user's message</h3>
            <div>
                {{ .Message }}
            </div>
            <h5>This is an automatic email. Do not reply</h5>
        </div>
</body>

</html>
