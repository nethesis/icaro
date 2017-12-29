import {Component} from '@angular/core';
import {Router} from '@angular/router';

@Component({
    moduleId:module.id,
    selector: 'error-page',
    templateUrl:'404.component.html',
    // styleUrls: ['404.component.css']
})

export class NotFoundComponent{
    constructor(public router: Router){}

}
