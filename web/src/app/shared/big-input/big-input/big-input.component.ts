import { Component, Input, ChangeDetectionStrategy } from '@angular/core';
import { NgClass } from '@angular/common';
import { MatCardModule } from '@angular/material/card';

@Component({
    selector: 'anms-big-input',
    templateUrl: './big-input.component.html',
    styleUrls: ['./big-input.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: true,
    imports: [MatCardModule, NgClass]
})
export class BigInputComponent {
  @Input()
  placeholder = '';

  @Input()
  value = '';

  @Input()
  disabled = false;

  hasFocus = false;
}
