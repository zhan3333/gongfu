import { enableProdMode, importProvidersFrom } from '@angular/core';


import { environment } from './environments/environment';
import { AppComponent } from './app/app/app.component';
import { SharedModule } from './app/shared/shared.module';
import { AppRoutingModule } from './app/app-routing.module';
import { CoreModule } from './app/core/core.module';
import { bootstrapApplication, BrowserModule } from '@angular/platform-browser';
import { provideAnimations } from '@angular/platform-browser/animations';
import localeZh from '@angular/common/locales/zh';
import { registerLocaleData } from '@angular/common';

registerLocaleData(localeZh)

if (environment.production) {
  enableProdMode();
}

bootstrapApplication(AppComponent, {
  providers: [
    importProvidersFrom(BrowserModule,
      // core
      CoreModule,
      // app
      AppRoutingModule, SharedModule),
    provideAnimations()
  ]
})
  .catch((err) => console.error(err));
