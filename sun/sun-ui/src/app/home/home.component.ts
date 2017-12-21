import {Component} from '@angular/core';
import {Router} from '@angular/router';


@Component({
	moduleId: module.id,
	selector: 'router-outlet',
	templateUrl:'home.component.html',
	styleUrls: ['home.component.css']
})


export class HomeComponent {
	constructor(public router: Router){}
}