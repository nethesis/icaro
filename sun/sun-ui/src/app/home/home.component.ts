import { LoginComponent } from './../login/login.component';
import {Component, } from '@angular/core';
import {Router} from '@angular/router';


@Component({
	moduleId: module.id,
	selector: 'home',
	
	templateUrl:'home.component.html',
	styleUrls: ['home.component.css']
})


export class HomeComponent {
	accountType:string;
	constructor(public router: Router){
		
	}
}