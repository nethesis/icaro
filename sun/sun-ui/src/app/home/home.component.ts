import { LoginComponent } from './../login/login.component';
import {Component, Input} from '@angular/core';
import {Router} from '@angular/router';
import {loginResponse} from './../_models/loginResponse'


@Component({
	moduleId: module.id,
	selector: 'home',
	
	templateUrl:'home.component.html',
	styleUrls: ['home.component.css']
})


export class HomeComponent {
	@Input('home') home: any;
	accountType:string;
	constructor(public router: Router){
	}
}