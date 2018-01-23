import { Component } from '@angular/core';
import { TranslateService } from 'ng2-translate';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
	
	 constructor(private translateService: TranslateService) {
	    this.translateService.addLangs(['en', 'it']);
      this.translateService.setDefaultLang('en');
      localStorage.setItem('language', this.translateService.getBrowserLang());
	    this.translateService.use(localStorage.getItem('language'));
  }
}
