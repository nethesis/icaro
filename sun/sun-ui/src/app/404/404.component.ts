import {Component} from '@angular/core';
import {Router} from '@angular/router';
import {Http, Headers} from '@angular/http'

@Component({
    moduleId:module.id,
    selector: 'my-app',
    templateUrl:'404.component.html',
    styleUrls: ['404.component.css']
})

export class NotFoundComponent{
    constructor (
        public router: Router, public http: Http
    ){}

}