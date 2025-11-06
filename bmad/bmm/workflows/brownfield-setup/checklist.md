# Brownfield Setup Workflow Validation Checklist

## Gap Detection
- [ ] All project components were correctly detected
- [ ] Missing components were accurately identified
- [ ] No false positives (marking existing as missing)
- [ ] No false negatives (missing actual gaps)

## Project Configuration
- [ ] Project type auto-detected or user-selected correctly
- [ ] BMAD module mapped appropriately
- [ ] Directory inference logic worked (Peps Ventures, Raegis Labs, XOGNOSIS)
- [ ] Inferred label and section are correct

## Component Installation (as applicable)
- [ ] Git initialized with appropriate .gitignore
- [ ] GitHub repository created and linked
- [ ] Linear project created with correct team and label
- [ ] WezTerm launcher updated in correct section
- [ ] BMAD v6 installed with appropriate module
- [ ] Foundation documents created (only missing ones)

## Integration
- [ ] GitHub repo linked to Linear project (if both created)
- [ ] All paths and URLs are valid
- [ ] No existing configuration was overwritten

## User Experience
- [ ] Clear feedback on what was found vs missing
- [ ] User asked for confirmation before applying changes
- [ ] Dry-run mode previewed changes correctly
- [ ] Skip flags respected user preferences
- [ ] Final summary accurately reflects changes made

## Documentation
- [ ] All created documents are complete and accurate
- [ ] No duplicate or conflicting documentation
- [ ] Project context properly captured

## Validation
- [ ] All newly created components function correctly
- [ ] No errors during workflow execution
- [ ] Project is in valid state after gap-filling
