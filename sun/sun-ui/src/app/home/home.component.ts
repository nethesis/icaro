import {Component} from '@angular/core';
import {Router} from '@angular/router';


@Component({
    moduleId: module.id,
    selector: 'r-outlet',
    templateUrl:'home.component.html'
})


export class HomeComponent {
    constructor(public router: Router){}
}