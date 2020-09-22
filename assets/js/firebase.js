import * as firebase from "firebase/app";
import "firebase/auth";
import * as firebaseui from "firebaseui";

function firebaseInit() {
    // Your web app's Firebase configuration
    var firebaseConfig = {
        apiKey: "AIzaSyCrgbSlSjbHLCFg0nkxVNHhCqkw4sMAAzo",
        authDomain: "gae-fire.firebaseapp.com",
        databaseURL: "https://gae-fire.firebaseio.com",
        projectId: "gae-fire",
        storageBucket: "gae-fire.appspot.com",
        messagingSenderId: "829754482529",
        appId: "1:829754482529:web:5fd191c7d6e0a49a706f20"
    };
    // Initialize Firebase
    firebase.initializeApp(firebaseConfig);
}

function firebaseInitUI() {
    var ui = new firebaseui.auth.AuthUI(firebase.auth());
    ui.start('#firebaseui-auth-container', {
        callbacks: {
            uiShown: function () {
                document.getElementById('loader').style.display = 'none';
            }
        },
        signInOptions: [
            firebase.auth.EmailAuthProvider.PROVIDER_ID
        ],
        signInFlow: 'popup',
        signInSuccessUrl: '/login_success',
    });
}

export {firebase, firebaseui, firebaseInitUI, firebaseInit}
