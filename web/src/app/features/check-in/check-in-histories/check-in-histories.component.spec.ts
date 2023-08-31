import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CheckInHistoriesComponent } from './check-in-histories.component';

describe('CheckInHistoriesComponent', () => {
  let component: CheckInHistoriesComponent;
  let fixture: ComponentFixture<CheckInHistoriesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    imports: [CheckInHistoriesComponent]
})
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CheckInHistoriesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
