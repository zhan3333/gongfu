import { ChangeDetectionStrategy, Component, Inject } from '@angular/core';
import { MAT_BOTTOM_SHEET_DATA, MatBottomSheetRef } from '@angular/material/bottom-sheet';
import { NgFor } from '@angular/common';
import { MatListModule } from '@angular/material/list';

@Component({
  selector: 'anms-bottom-sheet',
  template: `
    <mat-nav-list>
      <mat-list-item *ngFor="let k of data.keys()" (click)="select(k)">{{ data.get(k) }}</mat-list-item>
    </mat-nav-list>`,
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatListModule, NgFor],
})
export class BottomSheetComponent {
  constructor(
    @Inject(MAT_BOTTOM_SHEET_DATA) public data: Map<string, string>,
    private _bottomSheetRef: MatBottomSheetRef<BottomSheetComponent>
  ) {

  }

  select(k: string) {
    this._bottomSheetRef.dismiss(k)
  }
}
