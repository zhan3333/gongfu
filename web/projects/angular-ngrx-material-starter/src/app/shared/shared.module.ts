import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTabsModule as MatTabsModule } from '@angular/material/legacy-tabs';
import { MatLegacyInputModule as MatInputModule } from '@angular/material/legacy-input';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyCardModule as MatCardModule } from '@angular/material/legacy-card';
import { MatLegacyListModule as MatListModule } from '@angular/material/legacy-list';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatLegacySnackBarModule as MatSnackBarModule } from '@angular/material/legacy-snack-bar';
import { MatLegacySlideToggleModule as MatSlideToggleModule } from '@angular/material/legacy-slide-toggle';
import { MatDividerModule } from '@angular/material/divider';
import { MatLegacySliderModule as MatSliderModule } from '@angular/material/legacy-slider';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule, MatRippleModule } from '@angular/material/core';

import {
  FaIconLibrary,
  FontAwesomeModule
} from '@fortawesome/angular-fontawesome';
import {
  faAddressBook,
  faArrowLeft,
  faArrowRight,
  faBan,
  faBars,
  faBook,
  faBuilding,
  faCalendar,
  faCalendarAlt,
  faCaretDown,
  faCaretUp,
  faCheck,
  faEdit,
  faExclamationTriangle,
  faFile,
  faFilter,
  faLanguage,
  faLightbulb,
  faList,
  faPaintBrush,
  faPaperPlane,
  faPhone,
  faPlus,
  faSquare,
  faStream,
  faTasks,
  faTimes,
  faTrash,
  faUser,
  faUsers,
  faVideo,
  faWindowMaximize
} from '@fortawesome/free-solid-svg-icons';
import { faGithub, faMediumM } from '@fortawesome/free-brands-svg-icons';

import { BigInputComponent } from './big-input/big-input/big-input.component';
import { BigInputActionComponent } from './big-input/big-input-action/big-input-action.component';
import { RtlSupportDirective } from './rtl-support/rtl-support.directive';
import { CalendarHeaderComponent } from './calendar-header/calendar-header.component';
import { FormControlPipe } from './form-control.pipe';
import { faBolt } from '@fortawesome/free-solid-svg-icons/faBolt';
import { BottomSheetComponent } from './bottom-sheet.component';
import { MatBottomSheetModule } from '@angular/material/bottom-sheet';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,

    TranslateModule,

    MatButtonModule,
    MatSelectModule,
    MatTabsModule,
    MatInputModule,
    MatProgressSpinnerModule,
    MatChipsModule,
    MatCardModule,
    MatCheckboxModule,
    MatListModule,
    MatMenuModule,
    MatIconModule,
    MatTooltipModule,
    MatSnackBarModule,
    MatSlideToggleModule,
    MatDividerModule,
    MatDatepickerModule,
    MatNativeDateModule,

    FontAwesomeModule,
    MatBottomSheetModule,
    MatRippleModule
  ],
  declarations: [
    BigInputComponent,
    BigInputActionComponent,
    RtlSupportDirective,
    CalendarHeaderComponent,
    FormControlPipe,
    BottomSheetComponent
  ],
  exports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,

    TranslateModule,

    MatButtonModule,
    MatMenuModule,
    MatTabsModule,
    MatChipsModule,
    MatInputModule,
    MatProgressSpinnerModule,
    MatCheckboxModule,
    MatCardModule,
    MatListModule,
    MatSelectModule,
    MatIconModule,
    MatTooltipModule,
    MatSnackBarModule,
    MatSlideToggleModule,
    MatDividerModule,
    MatSliderModule,
    MatInputModule,
    MatDatepickerModule,
    MatNativeDateModule,

    FontAwesomeModule,

    BigInputComponent,
    BigInputActionComponent,
    RtlSupportDirective,
    FormControlPipe
  ]
})
export class SharedModule {
  constructor(faIconLibrary: FaIconLibrary) {
    faIconLibrary.addIcons(
      faGithub,
      faMediumM,
      faPlus,
      faEdit,
      faTrash,
      faTimes,
      faCaretUp,
      faCaretDown,
      faExclamationTriangle,
      faFilter,
      faTasks,
      faCheck,
      faSquare,
      faLanguage,
      faPaintBrush,
      faLightbulb,
      faWindowMaximize,
      faStream,
      faBook,
      faFile,
      faVideo,
      faArrowLeft,
      faArrowRight,
      faPaperPlane,
      faTrash,
      faCalendar,
      faCalendarAlt,
      faBolt,
      faPhone,
      faUser,
      faBuilding,
      faList,
      faAddressBook,
      faUsers,
      faBars,
      faBan
    );
  }
}
