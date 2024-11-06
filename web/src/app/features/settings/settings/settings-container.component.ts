import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { MatSelectChange, MatSelectModule } from '@angular/material/select';
import { MatSlideToggleChange, MatSlideToggleModule } from '@angular/material/slide-toggle';
import { select, Store } from '@ngrx/store';
import { Observable } from 'rxjs';

import { ROUTE_ANIMATIONS_ELEMENTS } from '../../../core/core.module';

import {
  actionSettingsChangeAnimationsElements,
  actionSettingsChangeAnimationsPage,
  actionSettingsChangeAutoNightMode,
  actionSettingsChangeLanguage,
  actionSettingsChangeStickyHeader,
  actionSettingsChangeTheme
} from '../../../core/settings/settings.actions';
import { SettingsState, State } from '../../../core/settings/settings.model';
import { selectSettings } from '../../../core/settings/settings.selectors';
import {
  faBars,
  faLanguage,
  faLightbulb,
  faPaintBrush,
  faStream,
  faWindowMaximize
} from '@fortawesome/free-solid-svg-icons';
import { TranslateModule } from '@ngx-translate/core';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatOptionModule } from '@angular/material/core';
import { FormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { MatIconModule } from '@angular/material/icon';
import { AsyncPipe, NgClass, NgFor, NgIf } from '@angular/common';
import { RtlSupportDirective } from '../../../shared/rtl-support/rtl-support.directive';
import { MatCardModule } from '@angular/material/card';

@Component({
  selector: 'anms-settings',
  templateUrl: './settings-container.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [RtlSupportDirective, NgIf, NgClass, MatIconModule, FontAwesomeModule, MatFormFieldModule, MatSelectModule, FormsModule, NgFor, MatOptionModule, MatSlideToggleModule, MatTooltipModule, AsyncPipe, TranslateModule, MatCardModule]
})
export class SettingsContainerComponent implements OnInit {
  public faLanguage = faLanguage;
  public faStream = faStream;
  public faBars = faBars;
  public faPaintBrush = faPaintBrush;
  public faLightbulb = faLightbulb;
  public faWindowMaximize = faWindowMaximize;
  routeAnimationsElements = ROUTE_ANIMATIONS_ELEMENTS;
  settings$: Observable<SettingsState> | undefined;

  themes = [
    {value: 'DEFAULT-THEME', label: 'green'},
    {value: 'BLUE-THEME', label: 'blue'},
    {value: 'LIGHT-THEME', label: 'light'},
    {value: 'NATURE-THEME', label: 'nature'},
    {value: 'BLACK-THEME', label: 'dark'}
  ];

  languages = [
    {value: 'en', label: 'English'},
    {value: 'de', label: 'Deutsch'},
    {value: 'sk', label: 'Slovenčina'},
    {value: 'fr', label: 'Français'},
    {value: 'es', label: 'Español'},
    {value: 'pt-br', label: 'Português'},
    {value: 'zh-cn', label: '简体中文'},
    {value: 'he', label: 'עברית'},
    {value: 'ar', label: 'اللغة العربية'}
  ];

  constructor(private store: Store<State>) {
  }

  ngOnInit() {
    this.settings$ = this.store.pipe(select(selectSettings));
  }

  onLanguageSelect(change: MatSelectChange) {
    this.store.dispatch(
      actionSettingsChangeLanguage({language: change.value})
    );
  }

  onThemeSelect(event: MatSelectChange) {
    this.store.dispatch(actionSettingsChangeTheme({theme: event.value}));
  }

  onAutoNightModeToggle(event: MatSlideToggleChange) {
    this.store.dispatch(
      actionSettingsChangeAutoNightMode({autoNightMode: event.checked})
    );
  }

  onStickyHeaderToggle(event: MatSlideToggleChange) {
    this.store.dispatch(
      actionSettingsChangeStickyHeader({stickyHeader: event.checked})
    );
  }

  onPageAnimationsToggle(event: MatSlideToggleChange) {
    this.store.dispatch(
      actionSettingsChangeAnimationsPage({pageAnimations: event.checked})
    );
  }

  onElementsAnimationsToggle(event: MatSlideToggleChange) {
    this.store.dispatch(
      actionSettingsChangeAnimationsElements({
        elementsAnimations: event.checked
      })
    );
  }
}
